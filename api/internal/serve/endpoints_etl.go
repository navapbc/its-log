package serve

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"regexp"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jadudm/its-log/internal/itslog"
// 	"github.com/jadudm/its-log/internal/sqlite"
// 	"github.com/jadudm/its-log/internal/sqlite/models"
// 	"github.com/spf13/viper"
// )

// type ETLParams struct {
// 	Date time.Time
// 	Name string
// }

// func addEtlEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *itslog.Event) {
// 	auth_adminV1 := rG.Group("/")
// 	permissions := []itslog.PermissionType{itslog.Admin, itslog.Test}
// 	auth_adminV1.Use(AuthMiddleWare(permissions))
// 	// Insert a new ETL step
// 	auth_adminV1.POST("etl/:date/:name", ETL)
// 	// Run an ETL step
// 	auth_adminV1.PUT("etl/:date/:name", ETL)
// 	// Retrieve the contents of a step, including the last run and run status
// 	auth_adminV1.GET("etl/:date/:name", ETL)
// 	// Remove a step
// 	auth_adminV1.DELETE("etl/:date/:name", ETL)
// 	// Combine a table from one DB into another DB
// 	// auth_adminV1.PUT("combine/:source/:destination/:table", Combine)
// }

// func (e *ETLParams) toSqliteFilename() string {
// 	return fmt.Sprintf("%s.sqlite", e.Date.Format(time.DateOnly))
// }

// func ETL(c *gin.Context) {
// 	date := c.Param("date")
// 	name := c.Param("name")

// 	// If things are malformed, return errors
// 	dateVal, err := time.Parse(time.DateOnly, date)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"method":  c.Request.Method,
// 			"message": fmt.Sprintf("%s is not YYYY-MM-DD", date),
// 		})
// 		return
// 	}

// 	// Bundle up params and call the correct method
// 	params := ETLParams{
// 		Date: dateVal,
// 		Name: name,
// 	}

// 	switch c.Request.Method {
// 	case http.MethodGet:
// 		get(c, params)
// 		return
// 	case http.MethodPost:
// 		post(c, params)
// 		return
// 	case http.MethodPut:
// 		put(c, params)
// 		return
// 	case http.MethodDelete:
// 		delete(c, params)
// 		return
// 	default:
// 		// It might not be possible to get here; Gin seems to
// 		// intercept unknown/underfined methods and return a 404.
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"method":  c.Request.Method,
// 			"message": "method not supported",
// 		})
// 		return
// 	}
// }

// type ETLPostBody struct {
// 	SQL string `json:"sql" binding:"required"`
// }

// type prePost struct {
// 	Pre  string
// 	Post string
// }

// // We have to trust the ETL pipeline. But, we can put some rudimentary
// // protections on it, so that data is not casually deleted, and direct
// // queries are not made against the events table---but instead through
// // the view appropriate to the key in question.
// // FIXME: If this is going to hang around, use a real parser.
// // https://github.com/rqlite/sql/blob/master/parser_test.go
// func filterEtlStatement(statement string, keyId string) error {
// 	// There's a trigger to prevent deletion, but we'll add that anyway.
// 	lower := strings.ToLower(statement)
// 	if strings.Contains(lower, "delete") {
// 		return errors.New("DELETE not permitted in ETL")
// 	}

// 	// Don't use the events table directly; that means any query that inclues
// 	// the table name with spaces around it
// 	prohibited_tables := []string{"itslog_events", "itslog_summary"}
// 	for _, table := range prohibited_tables {
// 		pre_posts := []prePost{{Pre: " ", Post: " "}, {Pre: " ", Post: ";"}, {Pre: " ", Post: ")"}}
// 		for _, pp := range pre_posts {
// 			if strings.Contains(lower, pp.Pre+table+pp.Post) {
// 				return errors.New("Direct use of the events or summary table is not permitted")
// 			}
// 		}
// 	}

// 	// If we see `itslog_events`, it must be followed by our _key.
// 	r, _ := regexp.Compile("itslog_events(_[a-zA-Z0-9]+)")
// 	matches := r.FindAllStringSubmatch(statement, -1)
// 	for _, m := range matches {
// 		if m[1] != "_"+keyId {
// 			return errors.New("Not allowed to reference itslog_events" + m[1])
// 		}
// 	}

// 	return nil
// }

// func replaceMagicStrings(s string, keyId string, params ETLParams) string {
// 	repl := map[string]string{
// 		"ITSLOG_EVENTS_TABLE": "itslog_events_" + keyId,
// 		"ITSLOG_KEY_ID":       keyId,
// 		"ITSLOG_DATE":         params.Date.Format("2006-01-02"),
// 	}

// 	query_string := string(s)
// 	for k, v := range repl {
// 		query_string = strings.ReplaceAll(query_string, k, v)
// 	}

// 	return query_string
// }

// // Create a restricted view for this key
// func createViewIfNotExists(db *sql.DB, keyId string) error {
// 	statement := `CREATE VIEW IF NOT EXISTS itslog_events_` + keyId
// 	statement += ` AS SELECT * FROM itslog_events WHERE key_id = `
// 	statement += fmt.Sprintf("'%s'", keyId) + ";"
// 	_, err := db.ExecContext(context.Background(), statement)
// 	if err != nil {
// 		return errors.New("could not create view for " + keyId)
// 	}
// 	return nil
// }

// // Insert a new ETL step
// func post(c *gin.Context, params ETLParams) {

// 	var body ETLPostBody
// 	// Call ShouldBindJSON to bind the incoming JSON to the newItem struct
// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		// If an error occurs (e.g., invalid JSON, missing required fields),
// 		// return a 400 Bad Request error
// 		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("should bind: %s", err.Error())})
// 		return
// 	}

// 	appId := itslog.GetOrPanic(c, "app_id")
// 	storage := &sqlite.SqliteStorage{
// 		Path:     viper.GetString("storage.path"),
// 		Filename: params.toSqliteFilename(),
// 		AppId:    appId,
// 	}
// 	storage.Init()
// 	defer storage.Close()

// 	// Replace magic placeholders
// 	magicked := replaceMagicStrings(body.SQL, keyId, params)
// 	log.Println("FILTERING SQL")
// 	log.Println(magicked)
// 	if e := filterEtlStatement(magicked, keyId); e != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  "error",
// 			"method":  c.Request.Method,
// 			"message": e.Error(),
// 			"date":    params.Date,
// 			"name":    params.Name,
// 		})
// 		return
// 	}

// 	if err := storage.GetQueries().InsertETL(context.Background(), models.InsertETLParams{
// 		KeyID: keyId,
// 		Name:  params.Name,
// 		Sql:   magicked,
// 	}); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"method":  c.Request.Method,
// 			"message": "could not save ETL step",
// 			"date":    params.Date,
// 			"name":    params.Name,
// 		})
// 		return
// 	}

// 	e := createViewIfNotExists(storage.GetDB(), keyId)
// 	if e != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"method":  c.Request.Method,
// 			"message": e.Error(),
// 			"date":    params.Date,
// 			"name":    params.Name,
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"status": "ok",
// 		"method": c.Request.Method,
// 		"date":   params.Date,
// 		"name":   params.Name,
// 	})
// }

// // Retrieve the contents of a step, including the last run and run status
// func get(c *gin.Context, params ETLParams) {

// 	keyId := GetKeyId(c)

// 	storage := &sqlite.SqliteStorage{
// 		Path:     viper.GetString("storage.path"),
// 		Filename: params.toSqliteFilename(),
// 		KeyId:    keyId,
// 	}
// 	storage.Init()
// 	defer storage.Close()

// 	row, err := storage.GetQueries().GetETL(context.Background(), models.GetETLParams{
// 		KeyID: keyId,
// 		Name:  params.Name,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status":  "error",
// 			"method":  c.Request.Method,
// 			"message": "could not find ETL step",
// 			"date":    params.Date,
// 			"name":    params.Name,
// 		})
// 		return
// 	}

// 	// It came back, so clean up the nullables, and send it back.
// 	last_run := ""
// 	if row.LastRun.Valid {
// 		last_run = row.LastRun.Time.Format("2006-01-02 15:04:05")
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":   "ok",
// 		"method":   c.Request.Method,
// 		"date":     params.Date,
// 		"name":     params.Name,
// 		"sql":      row.Sql,
// 		"last_run": last_run,
// 	})
// }

// // Run an ETL step
// func put(c *gin.Context, params ETLParams) {
// 	// Copypasta from above... :/
// 	keyId := GetKeyId(c)
// 	storage := &sqlite.SqliteStorage{
// 		Path:     viper.GetString("storage.path"),
// 		Filename: params.toSqliteFilename(),
// 		KeyId:    keyId,
// 	}
// 	storage.Init()
// 	defer storage.Close()

// 	// Get a Tx for making transaction requests.
// 	ctx := context.Background()
// 	tx, err := storage.GetDB().BeginTx(ctx, nil)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"method":  c.Request.Method,
// 			"message": "could not open transaction",
// 			"error":   err.Error(),
// 			"date":    params.Date,
// 			"name":    params.Name,
// 		})
// 		return
// 	}
// 	// Defer a rollback in case anything fails.
// 	defer tx.Rollback()

// 	qtx := storage.GetQueries().WithTx(tx)

// 	row, err := qtx.GetETL(context.Background(), models.GetETLParams{
// 		KeyID: keyId,
// 		Name:  params.Name,
// 	})

// 	if err != nil {
// 		msg := "could not find ETL step"
// 		log.Println(msg)
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status":  "error",
// 			"method":  c.Request.Method,
// 			"message": msg,
// 			"date":    params.Date,
// 			"name":    params.Name,
// 		})
// 		return
// 	}

// 	query_string := replaceMagicStrings(row.Sql, keyId, params)

// 	// Run the query
// 	_, err = tx.ExecContext(ctx, query_string)
// 	if err != nil {
// 		msg := "could not exec sql of ETL step"
// 		log.Println(msg)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"method":  c.Request.Method,
// 			"message": msg,
// 			"detail":  err.Error(),
// 			"date":    params.Date,
// 			"name":    params.Name,
// 		})
// 		return
// 	}

// 	if err = qtx.UpdateLastRun(ctx, models.UpdateLastRunParams{
// 		KeyID: keyId,
// 		Name:  params.Name,
// 	}); err != nil {
// 		msg := "could not update ETL metadata"
// 		log.Println(msg)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"method":  c.Request.Method,
// 			"message": msg,
// 			"date":    params.Date,
// 			"name":    params.Name,
// 		})
// 		return
// 	}

// 	// Commit the transaction.
// 	if err = tx.Commit(); err != nil {
// 		msg := "could not commit transaction"
// 		log.Println(msg)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"method":  c.Request.Method,
// 			"message": msg,
// 			"date":    params.Date,
// 			"name":    params.Name,
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"status": "ok",
// 		"method": c.Request.Method,
// 		"date":   params.Date,
// 		"name":   params.Name,
// 	})
// }

// func delete(c *gin.Context, params ETLParams) {

// 	c.JSON(http.StatusOK, gin.H{
// 		"status": "ok",
// 		"method": c.Request.Method,
// 		"date":   params.Date,
// 		"name":   params.Name,
// 	})
// }
