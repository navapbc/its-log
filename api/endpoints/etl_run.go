package endpoints

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	etl "github.com/navapbc/its-log/internal/base/etl/golang"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
)

// This is the handler if we call it directly from the API.
// We only get the gin context.
func RunEtl(c *gin.Context) {
	appId := base.GetOrPanic(c, "AppId")
	keyId := base.GetOrPanic(c, "KeyId")
	date := c.Param("date")
	name := c.Param("name")

	// ETL steps can have an arbitrary POST body?
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		body = []byte("{}")
	}
	payload := make(map[string]any)
	jsonErr := json.Unmarshal(body, &payload)
	if jsonErr != nil {
		payload = make(map[string]any)
	}

	s := types.NewStorage(appId)
	dateErr := s.SetDateYMD(date)
	if dateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not parse date; must be YYYY-MM-DD",
			"date":    date,
			"name":    name,
		})
		return
	}
	s.Init()

	pc, _, _, _ := runtime.Caller(0)
	funcName := runtime.FuncForPC(pc).Name()
	base.LoadDefaultEtlFiles(s, funcName)

	etlErr := runEtl(&types.RunEtlParams{
		AppId:   appId,
		KeyId:   keyId,
		GinCtx:  c,
		Storage: s,
		EtlName: name,
		Payload: payload,
	})

	if etlErr == nil {
		c.JSON(http.StatusCreated, gin.H{
			"status": "ok",
			"method": c.Request.Method,
			"date":   date,
			"name":   name,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "Golang ETL error: " + etlErr.Error(),
			"date":    date,
			"name":    name,
		})
		return
	}
}

func runEtl(etlP *types.RunEtlParams) error {
	row, err := etlP.Storage.Queries.GetETL(context.Background(), etlP.EtlName)

	if err != nil {
		msg := "could not find ETL step"
		log.Println(msg)
		return fmt.Errorf("%s: %s", msg, err.Error())
	}

	switch row.Kind {
	case "sql":
		err := etlRunSql(etlP, row, nil)
		if err != nil {
			return err
		}
	case "golang":
		err := etlRunGolang(etlP, row, nil)
		if err != nil {
			return err
		}
	case "starlark":
		err := etlRunStarlark(etlP, row, nil)
		if err != nil {
			return err
		}

	case "sequence":
		// The name of the step will have been loaded into the context.

	default:
		msg := "unknown kind of etl: " + row.Kind
		log.Println(msg)
		if etlP.GinCtx != nil {
			etlP.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  etlP.GinCtx.Request.Method,
				"message": msg,
				"date":    etlP.Storage.YYYYMMDD(),
				"name":    etlP.EtlName,
			})
		}
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	}

	updateErr := etlP.Storage.Queries.UpdateLastRun(context.Background(), etlP.EtlName)

	if updateErr != nil {
		msg := "could not update ETL metadata"
		log.Println(msg)
		if etlP.GinCtx != nil {
			etlP.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  etlP.GinCtx.Request.Method,
				"message": msg,
				"date":    etlP.Storage.YYYYMMDD(),
				"name":    etlP.EtlName,
			})
		}
		return fmt.Errorf("%s: %s", msg, err.Error())
	}

	return nil
}

func callWithParams(etlP *types.RunEtlParams, theSql string) *sql.Row {
	paramMap := map[string]string{
		"key_id": etlP.KeyId,
		"app_id": etlP.AppId,
		"date":   etlP.Storage.YYYYMMDD(),
	}

	args := make([]any, 0)
	for param, val := range paramMap {
		if strings.Contains(theSql, ":"+param) {
			// DEBUG LOG
			// log.Println(etlP.EtlName + " replacing :" + param + " with " + val)
			args = append(args, sql.Named(param, val))
		}
	}

	// An SQL step will read from and write to the DB.
	// Lock and release around this action.
	etlP.Storage.Lock()
	etlRow := etlP.Storage.GetDB().QueryRowContext(context.Background(), theSql, args...)
	etlP.Storage.Unlock()

	return etlRow
}

func etlRunSql(etlP *types.RunEtlParams, row models.GetETLRow, tx *sql.Tx) error {
	// Run the query
	if !row.Body.Valid {
		msg := "sql is null for ETL step"
		log.Println(msg)
		if etlP.GinCtx != nil {
			etlP.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  etlP.GinCtx.Request.Method,
				"message": msg,
				"detail":  msg,
				"date":    etlP.Storage.YYYYMMDD(),
				"name":    etlP.EtlName,
			})
		}
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	}

	// Expect all ETLs to return 0 or 1.
	etlRow := callWithParams(etlP, row.Body.String)

	var success int64 = -1
	err := etlRow.Scan(&success)
	switch err {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
	case nil:
		if success == constants.ETL_SUCCESS_VALUE {
			return nil
		} else {
			return fmt.Errorf("expected %d; %s returned %v", constants.ETL_SUCCESS_VALUE, row.Name, success)
		}
	default:
		msg := err.Error()
		if etlP.GinCtx != nil {
			etlP.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  etlP.GinCtx.Request.Method,
				"message": msg,
				"detail":  msg,
				"date":    etlP.Storage.YYYYMMDD(),
				"name":    etlP.EtlName,
			})
		}
		log.Println(msg)
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	}

	return nil
}

func etlRunGolang(etlP *types.RunEtlParams, row models.GetETLRow, tx *sql.Tx) error {
	functionName := row.Name
	if function, ok := etl.GolangETLMap[functionName]; ok {
		// Like the SQL actions, lock/unlock around what we do.
		// This saves the Golang ETLs from having to manage locks internal
		// to each individual action.
		etlP.Storage.Lock()
		err := function(etlP)
		etlP.Storage.Unlock()
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Function %s not found\n", functionName)
	}

	return nil
}

// -----------
// Running Starlark ETLs

func queryFun(etlP *types.RunEtlParams) func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var queryString string
		if err := starlark.UnpackArgs(
			"config", args, kwargs,
			"query", &queryString,
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
		}

		return resultList, nil

	}
}

func etlRunStarlark(etlP *types.RunEtlParams, row models.GetETLRow, tx *sql.Tx) error {
	// Run the query
	if !row.Body.Valid {
		msg := "sql is null for ETL step"
		log.Println(msg)
		if etlP.GinCtx != nil {
			etlP.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  etlP.GinCtx.Request.Method,
				"message": msg,
				"detail":  msg,
				"date":    etlP.Storage.YYYYMMDD(),
				"name":    etlP.EtlName,
			})
		}
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	}

	// Register a query function
	// This lets Starlark code call out and query the itslog_events table
	registrar := starlark.StringDict{"query": starlark.NewBuiltin("query", queryFun(etlP))}

	// The thread we'll execute the code in.
	thread := &starlark.Thread{Name: row.Name + "-starlark-thread"}

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
					var i int
					_ = starlark.AsInt(v, &i)
					srjs.Count = float64(i)

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
