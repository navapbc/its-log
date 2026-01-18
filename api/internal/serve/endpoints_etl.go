package serve

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/fsdb"
	"github.com/jadudm/its-log/internal/fsdb/models"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/spf13/viper"
)

type ETLParams struct {
	Date time.Time
	Name string
}

func addEtlEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *itslog.Event) {
	auth_adminV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Admin, itslog.Test}
	auth_adminV1.Use(AuthMiddleWare(permissions))
	// Insert a new ETL step
	auth_adminV1.POST("etl/:date/:name", ETL)
	// Run an ETL step
	auth_adminV1.PUT("etl/:date/:name", ETL)
	// Retrieve the contents of a step, including the last run and run status
	auth_adminV1.GET("etl/:date/:name", ETL)
	// Remove a step
	auth_adminV1.DELETE("etl/:date/:name", ETL)
	// Combine a table from one DB into another DB
	// auth_adminV1.PUT("combine/:source/:destination/:table", Combine)
}

func (e *ETLParams) formatDateTimeOnly() string {
	return fmt.Sprintf("%s.sqlite", e.Date.Format(time.DateOnly))
}

func ETL(c *gin.Context) {
	date := c.Param("date")
	name := c.Param("name")

	// If things are malformed, return errors
	dateVal, err := time.Parse(time.DateOnly, date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"method":  c.Request.Method,
			"message": fmt.Sprintf("%s is not YYYY-MM-DD", date),
		})
		return
	}

	// Bundle up params and call the correct method
	params := ETLParams{
		Date: dateVal,
		Name: name,
	}

	switch c.Request.Method {
	case http.MethodGet:
		get(c, params)
		return
	case http.MethodPost:
		post(c, params)
		return
	case http.MethodPut:
		put(c, params)
		return
	case http.MethodDelete:
		delete(c, params)
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
func post(c *gin.Context, params ETLParams) {

	var body ETLPostBody
	// Call ShouldBindJSON to bind the incoming JSON to the newItem struct
	if err := c.ShouldBindJSON(&body); err != nil {
		// If an error occurs (e.g., invalid JSON, missing required fields),
		// return a 400 Bad Request error
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("should bind: %s", err.Error())})
		return
	}

	appId := itslog.GetOrPanic(c, "app_id")
	storage := &fsdb.SqliteStorage{
		Path:     viper.GetString("storage.path"),
		Filename: params.formatDateTimeOnly(),
		AppId:    appId,
	}
	storage.Init()
	defer storage.Close()

	if err := storage.GetQueries().InsertETL(context.Background(), models.InsertETLParams{
		Name: params.Name,
		Sql:  body.SQL,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not save ETL step",
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"method": c.Request.Method,
		"date":   params.Date,
		"name":   params.Name,
	})
}

// Retrieve the contents of a step, including the last run and run status
func get(c *gin.Context, params ETLParams) {

	appId := itslog.GetOrPanic(c, itslog.ITSLOG_APPID)

	storage := &fsdb.SqliteStorage{
		Path:     viper.GetString("storage.path"),
		Filename: params.formatDateTimeOnly(),
		AppId:    appId,
	}
	storage.Init()
	defer storage.Close()

	row, err := storage.GetQueries().GetETL(context.Background(), params.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not find ETL step",
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}

	// It came back, so clean up the nullables, and send it back.
	last_run := ""
	if row.LastRun.Valid {
		last_run = row.LastRun.Time.Format("2006-01-02 15:04:05")
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"method":   c.Request.Method,
		"date":     params.Date,
		"name":     params.Name,
		"sql":      row.Sql,
		"last_run": last_run,
	})
}

// Run an ETL step
func put(c *gin.Context, params ETLParams) {
	// Copypasta from above... :/
	appId := itslog.GetOrPanic(c, itslog.ITSLOG_APPID)

	storage := &fsdb.SqliteStorage{
		Path:     viper.GetString("storage.path"),
		Filename: params.formatDateTimeOnly(),
		AppId:    appId,
	}
	storage.Init()
	defer storage.Close()

	// Get a Tx for making transaction requests.
	ctx := context.Background()
	tx, err := storage.GetDB().BeginTx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not open transaction",
			"error":   err.Error(),
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	qtx := storage.GetQueries().WithTx(tx)

	row, err := qtx.GetETL(context.Background(), params.Name)

	if err != nil {
		msg := "could not find ETL step"
		log.Println(msg)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": msg,
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}

	// Run the query
	_, err = tx.ExecContext(ctx, row.Sql)
	if err != nil {
		msg := "could not exec sql of ETL step"
		log.Println(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": msg,
			"detail":  err.Error(),
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}

	if err = qtx.UpdateLastRun(ctx, params.Name); err != nil {
		msg := "could not update ETL metadata"
		log.Println(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": msg,
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		msg := "could not commit transaction"
		log.Println(msg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": msg,
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"method": c.Request.Method,
		"date":   params.Date,
		"name":   params.Name,
	})
}

// FIXME: not implemented yet
func delete(c *gin.Context, params ETLParams) {

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"method": c.Request.Method,
		"date":   params.Date,
		"name":   params.Name,
	})
}
