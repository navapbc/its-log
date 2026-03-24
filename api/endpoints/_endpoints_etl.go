package serve

import (
	"context"

	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	etl "github.com/navapbc/its-log/internal/base/etl/golang"
	"github.com/navapbc/its-log/internal/base/models"
)

func addEtlEndpoints(rG *gin.RouterGroup) {
	auth_adminV1 := rG.Group("/")
	permissions := []base.PermissionType{base.Admin, base.Test}
	auth_adminV1.Use(AuthMiddleWare(permissions))

	// Insert a new ETL step
	auth_adminV1.POST("etl/:date", ETL)
	// Run an ETL step
	auth_adminV1.PUT("etl/:date/:name", ETL)
	// Retrieve the contents of a step, including the last run and run status
	auth_adminV1.GET("etl/:date/:name", ETL)
	// Combine a table from one DB into another DB
	// auth_adminV1.PUT("combine/:source/:destination/:table", Combine)
	auth_adminV1.GET("etl/reload/:date", ReloadEtl)
}

// ------------------------------------------------------------------------
// ETL
// ------------------------------------------------------------------------
func ETL(c *gin.Context) {
	sctx, err := base.NewServeCtx(c)
	if err != nil {
		return
	}
	defer sctx.Close()

	switch c.Request.Method {
	case http.MethodGet:
		get(sctx)
		return
	case http.MethodPost:
		post(sctx)
		return
	case http.MethodPut:
		put(sctx)
		return
	default:
		// It might not be possible to get here; Gin seems to
		// intercept unknown/underfined methods and return a 404.
		c.JSON(http.StatusBadRequest, gin.H{
			"method":  c.Request.Method,
			"message": "method not supported",
		})
		return
	}
}

type ETLPostBody struct {
	Name string `json:"name"`
	Kind string `json:"kind" binding:"required"`
	Body string `json:"body" binding:"required"`
}

// Insert a new ETL step
// Event godoc
// @Accept json
// @Description insert SQL for use as an ETL step
// @Failure 500	{object} base.Error
// @Param date YYYY-MM-DD formatted date for where to insert the step
// @Param name name of the ETL step (recommended: no spaces)
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /etl/{date}/{name} [post]
// @Schemes
// @Success 200 {object} base.Success
// @Summary insert SQL for use as an ETL step
// @Tags events
func post(sctx *base.ServeCtx) {

	var body ETLPostBody
	// Call ShouldBindJSON to bind the incoming JSON to the newItem struct
	if err := sctx.GinCtx.ShouldBindJSON(&body); err != nil {
		// If an error occurs (e.g., invalid JSON, missing required fields),
		// return a 400 Bad Request error
		sctx.GinCtx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("should bind: %s", err.Error())})
		return
	}

	if err := sctx.Storage.GetQueries().InsertETL(context.Background(), models.InsertETLParams{
		Name: sctx.Name,
		Kind: body.Kind,
		Body: sql.NullString{String: body.Body, Valid: true},
	}); err != nil {
		sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  sctx.RequestMethod,
			"message": "could not save ETL step",
			"date":    sctx.YYYYMMDD(),
			"name":    sctx.Name,
		})
		return
	}

	sctx.GinCtx.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"method": sctx.RequestMethod,
		"date":   sctx.YYYYMMDD(),
		"name":   sctx.Name,
	})
}

// Event godoc
// @Accept json
// @Description retrieve the contents of a step, including the last run and run status
// @Failure 500	{object} base.Error
// @Param date YYYY-MM-DD formatted date for where to insert the step
// @Param name name of the ETL step (recommended: no spaces)
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /etl/{date}/{name} [get]
// @Schemes
// @Success 200 {object} base.Success
// @Summary retrieve the contents of a step, including the last run and run status
// @Tags events
func get(sctx *base.ServeCtx) {
	err := reloadDefaultEtl(sctx)
	if err != nil {
		return
	}
	row, _ := possiblyFetchEtlRows(sctx)

	// It came back, so clean up the nullables, and send it back.
	last_run := ""
	if row.LastRun.Valid {
		last_run = row.LastRun.Time.Format("2006-01-02 15:04:05")
	}
	sctx.GinCtx.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"method":   sctx.RequestMethod,
		"date":     sctx.YYYYMMDD(),
		"name":     sctx.Name,
		"sql":      row.Body,
		"last_run": last_run,
	})
}

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

// Event godoc
// @Accept json
// @Description run an ETL step
// @Failure 500	{object} base.Error
// @Param date YYYY-MM-DD formatted date for where to insert the step
// @Param name name of the ETL step (recommended: no spaces)
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /etl/{date}/{name} [get]
// @Schemes
// @Success 200 {object} base.Success
// @Summary run an ETL step
// @Tags events
func put(sctx *base.ServeCtx) error {
	// This returns an error because it gets used to run sequences of steps, too.
	// In that context, the gin response is used, *and* the error.

	isSequence, _ := sctx.GinCtx.Get("isSequence")
	useContext := true
	if isSequence != nil {
		useContext = !isSequence.(bool)
	}
	sctx.GinCtx.Set("useContext", useContext)

	// Get a Tx for making transaction requests.
	// ctx := context.Background()
	// tx, err := sctx.Storage.GetDB().BeginTx(ctx, nil)
	// if err != nil {
	// 	msg := "could not open transaction"
	// 	log.Println(msg)
	// 	if useContext {
	// 		sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
	// 			"status":  "error",
	// 			"method":  sctx.RequestMethod,
	// 			"message": msg,
	// 			"error":   err.Error(),
	// 			"date":    sctx.YYYYMMDD(),
	// 			"name":    sctx.Name,
	// 		})
	// 	}
	// 	return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	// }
	// // Defer a rollback in case anything fails.
	// defer tx.Rollback()

	//qtx := sctx.Storage.GetQueries().WithTx(tx)

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

func ReloadEtl(c *gin.Context) {
	sctx, err := base.NewServeCtx(c)
	if err != nil {
		return
	}
	defer sctx.Close()
	reloadDefaultEtl(sctx)
}

func possiblyFetchEtlRows(sctx *base.ServeCtx) (*models.GetETLRow, bool) {
	row, err := sctx.Storage.GetQueries().GetETL(context.Background(), "sentinel")
	if err != nil {
		return nil, false
	}
	return &row, true
}

func reloadDefaultEtl(sctx *base.ServeCtx) error {
	// This is funky. It also happens on the first attempt to init() the DB.
	// But, that happens on buffer flush.
	// And, we don't know when someone is logging for (e.g. it could be in the past, if this
	// is for testing.) So, we check here to make sure the ETL has been loaded.
	_, exist := possiblyFetchEtlRows(sctx)
	if !exist {
		sctx.Storage.LoadDefaultEtlFiles()
		_, exist = possiblyFetchEtlRows(sctx)
	}
	if !exist {
		sctx.GinCtx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"method":  sctx.RequestMethod,
			"message": "could not find ETL step",
			"date":    sctx.YYYYMMDD(),
			"name":    sctx.Name,
		})
		return errors.New("could not find ETL step")
	}
	return nil
}
