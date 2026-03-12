package serve

import (
	"bufio"
	"context"

	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/fsdb/models"
	"github.com/jadudm/its-log/internal/itslog"
)

func addEtlEndpoints(rG *gin.RouterGroup) {
	auth_adminV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Admin, itslog.Test}
	auth_adminV1.Use(AuthMiddleWare(permissions))

	// Insert a new ETL step
	auth_adminV1.POST("etl/:date/:name", ETL)
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
	sctx, err := NewServeCtx(c)
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
	SQL string `json:"sql" binding:"required"`
}

// Insert a new ETL step
// Event godoc
// @Accept json
// @Description insert SQL for use as an ETL step
// @Failure 500	{object} itslog.Error
// @Param date YYYY-MM-DD formatted date for where to insert the step
// @Param name name of the ETL step (recommended: no spaces)
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /etl/{date}/{name} [post]
// @Schemes
// @Success 200 {object} itslog.Success
// @Summary insert SQL for use as an ETL step
// @Tags events
func post(sctx *ServeCtx) {

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
		Sql:  body.SQL,
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
// @Failure 500	{object} itslog.Error
// @Param date YYYY-MM-DD formatted date for where to insert the step
// @Param name name of the ETL step (recommended: no spaces)
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /etl/{date}/{name} [get]
// @Schemes
// @Success 200 {object} itslog.Success
// @Summary retrieve the contents of a step, including the last run and run status
// @Tags events
func get(sctx *ServeCtx) {
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
		"sql":      row.Sql,
		"last_run": last_run,
	})
}

func getRequestedParamKeys(sql string) ([]string, bool) {
	// Insert params
	// Valid params include
	// -- params: key_id app_id
	results := make([]string, 0)
	has_params := false

	scanner := bufio.NewScanner(strings.NewReader(sql))

	var dropped_leading_parts []string

	for scanner.Scan() {
		line := scanner.Text()
		// Make sure the line starts with -- params:
		is_params := strings.HasPrefix(line, "-- params:")
		pattern := regexp.MustCompile(`\S+(.*?)`)
		if is_params {
			results = pattern.FindAllString(line, -1)
			has_params = true
			dropped_leading_parts = results[2:]
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return dropped_leading_parts, has_params
}

// Event godoc
// @Accept json
// @Description run an ETL step
// @Failure 500	{object} itslog.Error
// @Param date YYYY-MM-DD formatted date for where to insert the step
// @Param name name of the ETL step (recommended: no spaces)
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /etl/{date}/{name} [get]
// @Schemes
// @Success 200 {object} itslog.Success
// @Summary run an ETL step
// @Tags events
func put(sctx *ServeCtx) error {
	// This returns an error because it gets used to run sequences of steps, too.
	// In that context, the gin response is used, *and* the error.

	isSequence, _ := sctx.GinCtx.Get("isSequence")
	useContext := true
	if isSequence != nil {
		useContext = !isSequence.(bool)
	}

	// Get a Tx for making transaction requests.
	ctx := context.Background()
	tx, err := sctx.Storage.GetDB().BeginTx(ctx, nil)
	if err != nil {
		msg := "could not open transaction"
		log.Println(msg)
		if useContext {
			sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  sctx.RequestMethod,
				"message": msg,
				"error":   err.Error(),
				"date":    sctx.YYYYMMDD(),
				"name":    sctx.Name,
			})
		}
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	qtx := sctx.Storage.GetQueries().WithTx(tx)

	row, err := qtx.GetETL(context.Background(), sctx.Name)

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
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusNotFound))
	}

	// We allow queries to specity a list of parameters of the form
	// -- params: a b c
	// where a, b, c must be in the set defined by the map below
	// keys, has_params := getRequestedParamKeys(row.Sql)
	// mapped_params := make(map[string]string, 0)
	// mapped := make([]any, 0)
	// if has_params {
	// 	mapped_params["key_id"] = sctx.KeyId
	// 	mapped_params["app_id"] = sctx.AppId
	// 	mapped_params["date"] = sctx.YYYYMMDD()
	// 	for _, p := range keys {
	// 		if v, ok := mapped_params[p]; ok {
	// 			mapped = append(mapped, v)
	// 		}
	// 	}
	// }

	named := make([]sql.NamedArg, 0)
	named = append(named, sql.Named("key_id", sctx.KeyId))
	named = append(named, sql.Named("app_id", sctx.AppId))
	named = append(named, sql.Named("date", sctx.YYYYMMDD()))
	// Run the query

	_, err = tx.ExecContext(ctx, row.Sql,
		sql.Named("key_id", sctx.KeyId),
		sql.Named("app_id", sctx.AppId),
		sql.Named("date", sctx.YYYYMMDD()))

	if err != nil {
		msg := "could not exec sql of ETL step"
		log.Println(msg)
		if useContext {
			sctx.GinCtx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"method":  sctx.RequestMethod,
				"message": msg,
				"detail":  err.Error(),
				"date":    sctx.YYYYMMDD(),
				"name":    sctx.Name,
			})
		}
		return fmt.Errorf("%s: %s", msg, http.StatusText(http.StatusInternalServerError))
	}

	if err = qtx.UpdateLastRun(ctx, sctx.Name); err != nil {
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
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		msg := "could not commit transaction"
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
	sctx, err := NewServeCtx(c)
	if err != nil {
		return
	}
	defer sctx.Close()
	reloadDefaultEtl(sctx)
}

func possiblyFetchEtlRows(sctx *ServeCtx) (*models.GetETLRow, bool) {
	row, err := sctx.Storage.GetQueries().GetETL(context.Background(), "sentinel")
	if err != nil {
		return nil, false
	}
	return &row, true
}

func reloadDefaultEtl(sctx *ServeCtx) error {
	// This is funky. It also happens on the first attempt to init() the DB.
	// But, that happens on buffer flush.
	// And, we don't know when someone is logging for (e.g. it could be in the past, if this
	// is for testing.) So, we check here to make sure the ETL has been loaded.
	_, exist := possiblyFetchEtlRows(sctx)
	if !exist {
		sctx.Storage.LoadDefaultEtlSql()
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
