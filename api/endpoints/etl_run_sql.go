package endpoints

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	etl "github.com/navapbc/its-log/internal/base/etl/golang"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
)

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
			ser := types.NewStandardErrorResponse(etlP, errors.New(msg))
			ser.SetStatus(http.StatusInternalServerError).Send(msg)
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
			ser := types.NewStandardErrorResponse(etlP, err)
			ser.SetStatus(http.StatusInternalServerError).Send(msg)
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
