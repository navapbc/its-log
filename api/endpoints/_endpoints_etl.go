package serve

import (
	"context"

	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/schema/models"
)

func etlRunSql(sctx *base.ServeCtx, row models.GetETLRow, tx *sql.Tx) error {
	useContext, _ := sctx.GinCtx.Get("useContext")

	named := make([]sql.NamedArg, 0)
	named = append(named, sql.Named("key_id", sctx.KeyId))
	named = append(named, sql.Named("app_id", sctx.AppId))
	named = append(named, sql.Named("date", sctx.YYYYMMDD()))

	// Run the query
	if !row.Body.Valid {
		msg := "sql is null for ETL step"
		log.Println(msg)
		if useContext.(bool) {
			sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  sctx.RequestMethod,
				"message": msg,
				"detail":  msg,
				"date":    sctx.YYYYMMDD(),
				"name":    sctx.Name,
			})
		}
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	}

	// Expect all ETLs to return 0 or 1.
	etlRow := sctx.Storage.GetDB().QueryRowContext(context.Background(), row.Body.String,
		sql.Named("key_id", sctx.KeyId),
		sql.Named("app_id", sctx.AppId),
		sql.Named("date", sctx.YYYYMMDD()))

	var success int64
	switch err := etlRow.Scan(&success); err {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
	case nil:
		if success == 1 {
			return nil
		} else {
			return fmt.Errorf("expected 1; %s returned %v", row.Name, success)
		}
	default:
		panic(err)
	}

	return nil
}

func etlRunGolang(sctx *base.ServeCtx, row models.GetETLRow, tx *sql.Tx) error {
	functionName := row.Name
	if function, ok := etl.GolangETLMap[functionName]; ok {
		function(sctx, tx)
	} else {
		fmt.Printf("Function %s not found\n", functionName)
	}

	return nil
}

func put(sctx *base.ServeCtx) error {
	// This returns an error because it gets used to run sequences of steps, too.
	// In that context, the gin response is used, *and* the error.

	isSequence, _ := sctx.GinCtx.Get("isSequence")
	useContext := true
	if isSequence != nil {
		useContext = !isSequence.(bool)
	}
	sctx.GinCtx.Set("useContext", useContext)

	row, err := sctx.Storage.GetQueries().GetETL(context.Background(), sctx.Name)

	if err != nil {
		msg := "could not find ETL step"
		log.Println(msg)
		if useContext {
			sctx.GinCtx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"method":  sctx.RequestMethod,
				"message": msg,
				"date":    sctx.YYYYMMDD(),
				"name":    sctx.Name,
			})
		}
		return fmt.Errorf("%s: %s", msg, err.Error())
	}

	switch row.Kind {
	case "sql":
		err := etlRunSql(sctx, row, nil)
		if err != nil {
			return err
		}
	case "golang":
		err := etlRunGolang(sctx, row, nil)
		if err != nil {
			return err
		}
	case "sequence":
		// The name of the step will have been loaded into the context.

	default:
		msg := "unknown kind of etl: " + row.Kind
		log.Println(msg)
		if useContext {
			sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  sctx.RequestMethod,
				"message": msg,
				"date":    sctx.YYYYMMDD(),
				"name":    sctx.Name,
			})
		}
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	}

	err = sctx.Storage.GetQueries().UpdateLastRun(context.Background(), sctx.Name)

	if err != nil {
		msg := "could not update ETL metadata"
		log.Println(msg)
		if useContext {
			sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  sctx.RequestMethod,
				"message": msg,
				"date":    sctx.YYYYMMDD(),
				"name":    sctx.Name,
			})
		}
		return fmt.Errorf("%s: %s", msg, err.Error())
	}

	// Commit the transaction.
	// if err = tx.Commit(); err != nil {
	// 	msg := "could not commit transaction"
	// 	log.Println(msg)
	// 	if useContext {
	// 		sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
	// 			"status":  "error",
	// 			"method":  sctx.RequestMethod,
	// 			"message": msg,
	// 			"date":    sctx.YYYYMMDD(),
	// 			"name":    sctx.Name,
	// 		})
	// 	}
	// 	return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	// }

	if useContext {
		sctx.GinCtx.JSON(http.StatusCreated, gin.H{
			"status": "ok",
			"method": sctx.RequestMethod,
			"date":   sctx.YYYYMMDD(),
			"name":   sctx.Name,
		})
	}

	return nil
}

// func ReloadEtl(c *gin.Context) {
// 	sctx, err := base.NewServeCtx(c)
// 	if err != nil {
// 		return
// 	}
// 	defer sctx.Close()
// 	reloadDefaultEtl(sctx)
// }

// func possiblyFetchEtlRows(sctx *base.ServeCtx) (*models.GetETLRow, bool) {
// 	row, err := sctx.Storage.GetQueries().GetETL(context.Background(), "sentinel")
// 	if err != nil {
// 		return nil, false
// 	}
// 	return &row, true
// }

// func reloadDefaultEtl(sctx *base.ServeCtx) error {
// 	// This is funky. It also happens on the first attempt to init() the DB.
// 	// But, that happens on buffer flush.
// 	// And, we don't know when someone is logging for (e.g. it could be in the past, if this
// 	// is for testing.) So, we check here to make sure the ETL has been loaded.
// 	_, exist := possiblyFetchEtlRows(sctx)
// 	if !exist {
// 		sctx.Storage.LoadDefaultEtlFiles()
// 		_, exist = possiblyFetchEtlRows(sctx)
// 	}
// 	if !exist {
// 		sctx.GinCtx.JSON(http.StatusNotFound, gin.H{
// 			"status":  "error",
// 			"method":  sctx.RequestMethod,
// 			"message": "could not find ETL step",
// 			"date":    sctx.YYYYMMDD(),
// 			"name":    sctx.Name,
// 		})
// 		return errors.New("could not find ETL step")
// 	}
// 	return nil
// }
