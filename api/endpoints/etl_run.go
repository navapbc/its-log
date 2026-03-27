package endpoints

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	etl "github.com/navapbc/its-log/internal/base/etl/golang"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
)

// This is the handler if we call it directly from the API.
// We only get the gin context.
func RunEtlHandler(c *gin.Context) {
	appId := base.GetOrPanic(c, "AppId")
	keyId := base.GetOrPanic(c, "KeyId")
	date := c.GetString("date")
	name := c.GetString("name")

	s := types.NewStorage(appId)
	err := s.SetDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not parse date; must be YYYY-MM-DD",
			"date":    date,
			"name":    name,
		})
		return
	}

	err = runEtl(&types.RunEtlParams{
		AppId:   appId,
		KeyId:   keyId,
		GinCtx:  c,
		Storage: s,
		EtlName: name,
	})

	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"status": "ok",
			"method": c.Request.Method,
			"date":   date,
			"name":   name,
		})
	}
}

func runEtl(etlP *types.RunEtlParams) error {
	row, err := etlP.Storage.Queries.GetETL(context.Background(), etlP.EtlName)

	if err != nil {
		msg := "could not find ETL step"
		log.Println(msg)
		if etlP.GinCtx != nil {
			etlP.GinCtx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"method":  etlP.GinCtx.Request.Method,
				"message": msg,
				"date":    etlP.Storage.YYYYMMDD(),
				"name":    etlP.EtlName,
			})
		}
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

	err = etlP.Storage.Queries.UpdateLastRun(context.Background(), etlP.EtlName)

	if err != nil {
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
			args = append(args, sql.Named(param, val))
		}
	}
	// log.Println(theSql)
	etlRow := etlP.Storage.GetDB().QueryRowContext(context.Background(), theSql, args...)
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
		log.Println(etlP)
		panic(err)
	}

	return nil
}

func etlRunGolang(etlP *types.RunEtlParams, row models.GetETLRow, tx *sql.Tx) error {
	functionName := row.Name
	if function, ok := etl.GolangETLMap[functionName]; ok {
		err := function(etlP)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Function %s not found\n", functionName)
	}

	return nil
}
