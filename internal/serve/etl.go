package serve

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/sqlite"
	"github.com/jadudm/its-log/internal/sqlite/models"
	"github.com/spf13/viper"
)

type ETLParams struct {
	Date time.Time
	Name string
}

func (e *ETLParams) toSqliteFilename() string {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storage := &sqlite.SqliteStorage{
		Path:     viper.GetString("storage.path"),
		Filename: params.toSqliteFilename(),
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

	storage := &sqlite.SqliteStorage{
		Path:     viper.GetString("storage.path"),
		Filename: params.toSqliteFilename(),
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
	storage := &sqlite.SqliteStorage{
		Path:     viper.GetString("storage.path"),
		Filename: params.toSqliteFilename(),
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
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not find ETL step",
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}

	// Run the query
	if _, err := tx.ExecContext(ctx, string(row.Sql)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not run ETL step",
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}

	if err = qtx.UpdateLastRun(ctx, params.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not update ETL metadata",
			"date":    params.Date,
			"name":    params.Name,
		})
		return
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not commit transaction",
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

func delete(c *gin.Context, params ETLParams) {

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"method": c.Request.Method,
		"date":   params.Date,
		"name":   params.Name,
	})
}
