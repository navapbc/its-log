package endpoints

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
	"go.starlark.net/lib/json"
	"go.starlark.net/lib/math"
	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

// -----------
// Running Starlark ETLs

func queryFun(etlP *types.RunEtlParams) func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var queryTable string
		var queryString string
		if err := starlark.UnpackArgs(
			"config", args, kwargs,
			"table", &queryTable, "query", &queryString,
		); err != nil {
			return starlark.None, fmt.Errorf("config: %s", err)
		}

		paramMap := map[string]string{
			"key_id": etlP.KeyId,
			"app_id": etlP.AppId,
			"date":   etlP.Storage.YYYYMMDD(),
		}

		sqlArgs := make([]any, 0)
		for param, val := range paramMap {
			if strings.Contains(queryString, ":"+param) {
				log.Println(etlP.EtlName + " replacing :" + param + " with " + val)
				sqlArgs = append(sqlArgs, sql.Named(param, val))
			}
		}

		// An SQL step will read from and write to the DB.
		// Lock and release around this action.
		rows, err := etlP.Storage.GetDB().QueryContext(context.Background(), queryString, sqlArgs...)
		if err != nil {
			return starlark.None, err
		}

		// DEBUG LOG
		// log.Println("golang", queryString)

		resultList := starlark.NewList([]starlark.Value{})
		for rows.Next() {
			switch queryTable {
			case "events":
				var row types.EventRow
				if err := rows.Scan(&row.ID, &row.Timestamp, &row.KeyId, &row.Cluster, &row.Tags, &row.Value); err != nil {
					return starlark.None, err
				} else {
					d := starlark.NewDict(10)
					d.SetKey(starlark.String("id"), starlark.String(row.ID))
					d.SetKey(starlark.String("timestamp"), starlark.String(row.Timestamp))
					d.SetKey(starlark.String("key_id"), starlark.String(row.KeyId))
					d.SetKey(starlark.String("cluster"), starlark.String(row.Cluster.String))
					d.SetKey(starlark.String("tags"), starlark.String(row.Tags))
					d.SetKey(starlark.String("value"), starlark.String(row.Value.String))
					resultList.Append(d)
				}
			case "summary":
				var row types.SummaryRow
				if err := rows.Scan(&row.ID, &row.LastRun, &row.Date, &row.KeyId, &row.Operation, &row.Tags, &row.Value, &row.Count, &row.Hash); err != nil {
					return starlark.None, err
				} else {
					d := starlark.NewDict(10)
					d.SetKey(starlark.String("id"), starlark.String(row.ID))
					d.SetKey(starlark.String("last_run"), starlark.MakeInt64(row.LastRun))
					d.SetKey(starlark.String("date"), starlark.String(row.KeyId))
					d.SetKey(starlark.String("key_id"), starlark.String(row.KeyId))
					d.SetKey(starlark.String("operation"), starlark.String(row.Operation))
					d.SetKey(starlark.String("tags"), starlark.String(row.Tags.String))
					d.SetKey(starlark.String("value"), starlark.String(row.Value.String))
					d.SetKey(starlark.String("count"), starlark.Float(row.Count))
					d.SetKey(starlark.String("hash"), starlark.String(row.Hash.String))
					resultList.Append(d)
				}
			}
		}

		return resultList, nil

	}
}

func etlRunStarlark(etlP *types.RunEtlParams, row models.GetETLRow, tx *sql.Tx) error {
	// Run the query
	if !row.Body.Valid {
		msg := "sql is null for ETL step"
		err := fmt.Errorf("etl body is not valid")
		ser := types.NewStandardErrorResponse(etlP, err)
		ser.SetStatus(http.StatusInternalServerError).Send(msg)
		return err
	}

	// The thread we'll execute the code in.
	thread := &starlark.Thread{Name: row.Name + "-starlark-thread"}
	// This loads standard modules, including JSON manipulation.
	// See
	// https://github.com/google/starlark-go/blob/fadfc96def35ea95e7f2bd9952256d4db1d80d91/cmd/starlark/starlark.go#L98
	starlark.Universe["json"] = json.Module
	starlark.Universe["time"] = time.Module
	starlark.Universe["math"] = math.Module

	// Register a query function
	// This lets Starlark code call out and query the itslog_events table
	registrar := starlark.StringDict{"query": starlark.NewBuiltin("query", queryFun(etlP))}

	// This evals the file, which does not *yet* execute the code.
	// Note we register "query" as a function in the interpreted namespace.
	globals, err := starlark.ExecFileOptions(&syntax.FileOptions{Recursion: false}, thread, "", row.Body.String, registrar)
	if err != nil {
		return err
	}

	// Call Starlark function from Go.
	// Specifically, we're calling "summarize", which must be provided.
	fn, ok := globals["summarize"].(*starlark.Function)
	if !ok {
		log.Printf("value %v is not a Starlark function", globals["summarize"])
		return fmt.Errorf("value %v is not a Starlark function", globals["summarize"])
	}

	// Here is the call.
	v, err := starlark.Call(thread, fn, nil, nil)
	if err != nil {
		return err
	}

	// Now, we're going to construct a list of rows that are suitable for insertion into the
	// summary table from this result. We should get back a Starlark list, which will
	// contain starlark dictionaries. These need to be turned into Golang structs.
	// FIXME: what if someone wants `summary` rows? This does not work.
	// consider using map[string]any instead.
	arr := make([]models.ItslogSummary, 0)
	for elem := range v.(*starlark.List).Elements() {
		srjs := models.ItslogSummary{}
		for _, k := range elem.(*starlark.Dict).Keys() {
			s, _ := starlark.AsString(k)
			switch s {
			case "operation":
				v, found, _ := elem.(*starlark.Dict).Get(k)
				if found {
					srjs.Operation, _ = starlark.AsString(v)
				}
			case "tags":
				v, found, _ := elem.(*starlark.Dict).Get(k)
				if found {
					srjs.Tags, _ = starlark.AsString(v)
				}
			case "value":
				v, found, _ := elem.(*starlark.Dict).Get(k)
				if found {
					srjs.Tags, _ = starlark.AsString(v)
				}
			case "count":
				v, found, _ := elem.(*starlark.Dict).Get(k)
				if found {
					var f float64
					f, _ = starlark.AsFloat(v)
					srjs.Count = float64(f)

				}
			}
		}
		arr = append(arr, srjs)
	}

	// Now, write them to the summary table.
	etlP.Storage.Lock()
	for _, e := range arr {
		summ := models.InsertSummaryParams{
			KeyID:     etlP.KeyId,
			Date:      etlP.Storage.ILTime.AsYYYYMMDD(),
			Operation: e.Operation,
			Tags:      e.Tags,
			Value:     e.Value,
			Count:     float64(e.Count),
		}
		summ.Hash = sql.NullString{String: e.ReturnHash(), Valid: true}

		err := etlP.Storage.Queries.InsertSummary(context.Background(), summ)
		if err != nil {
			etlP.Storage.Unlock()
			return err
		}
	}
	etlP.Storage.Unlock()

	return nil
}
