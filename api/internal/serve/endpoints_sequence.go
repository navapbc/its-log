package serve

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/fsdb"
	"github.com/jadudm/its-log/internal/fsdb/models"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/spf13/viper"
)

func addSequenceEndpoints(rG *gin.RouterGroup) {
	auth_adminV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Admin, itslog.Test}
	auth_adminV1.Use(AuthMiddleWare(permissions))

	// Add a sequence
	auth_adminV1.POST("sequence", Sequence)
	// Run a sequence
	auth_adminV1.GET("sequence/:date/:name", Sequence)
	// Delete a sequence
	auth_adminV1.DELETE("sequence", Sequence)

}

// ------------------------------------------------------------------------
// Sequence
// ------------------------------------------------------------------------

type SequencePostBody struct {
	Name  string   `json:"name" binding:"required"`
	Date  string   `json:"date" binding:"required"`
	Steps []string `json:"steps" binding:"required"`
}

func Sequence(c *gin.Context) {
	// Bundle up params and call the correct method

	switch c.Request.Method {
	case http.MethodPost:
		var body SequencePostBody
		// Call ShouldBindJSON to bind the incoming JSON to the newItem struct
		if err := c.ShouldBindJSON(&body); err != nil {
			// If an error occurs (e.g., invalid JSON, missing required fields),
			// return a 400 Bad Request error
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("should bind: %s", err.Error())})
			return
		}

		// If things are malformed, return errors
		_, err := time.Parse(time.DateOnly, body.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"method":  c.Request.Method,
				"message": fmt.Sprintf("%s is not YYYY-MM-DD", body.Date),
			})
			return
		}
		seq_post(c, body)
		return
	case http.MethodGet:
		seq_get(c)
		return
	case http.MethodDelete:
		// Bundle up params and call the correct method
		// seq_delete(c)
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

func seq_post(c *gin.Context, seqPostbody SequencePostBody) {

	appId := itslog.GetOrPanic(c, itslog.ITSLOG_APPID)
	keyId := itslog.GetOrPanic(c, itslog.ITSLOG_KEYID)

	storage := &fsdb.SqliteStorage{
		AppId: appId,
		Date:  seqPostbody.Date,
		Kind:  fsdb.NamedDatabase,
		Path:  viper.GetString("storage.path"),
	}

	storage.Init()
	defer storage.Close()

	if err := storage.GetQueries().InsertSequence(context.Background(), models.InsertSequenceParams{
		KeyID: keyId,
		Name:  seqPostbody.Name,
		Steps: strings.Join(seqPostbody.Steps, ","),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not save ETL step",
			"date":    seqPostbody.Date,
			"name":    seqPostbody.Name,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"method": c.Request.Method,
		"date":   seqPostbody.Date,
		"name":   seqPostbody.Name,
	})

}

// Run a sequence immediately, synchronously.
func seq_get(c *gin.Context) {
	name := c.Param("name")
	date := c.Param("date")

	dateVal, _ := time.Parse(time.DateOnly, date)
	var err error = nil
	errStep := ""

	appId := itslog.GetOrPanic(c, itslog.ITSLOG_APPID)

	storage := &fsdb.SqliteStorage{
		AppId: appId,
		Date:  date,
		Kind:  fsdb.NamedDatabase,
		Path:  viper.GetString("storage.path"),
	}

	storage.Init()
	defer storage.Close()

	seq, err := storage.GetQueries().GetSequence(context.Background(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not find sequence",
			"date":    date,
			"name":    name,
		})
		return
	}

	steps := strings.Split(seq, ",")

	for _, step := range steps {
		// Bundle up params and call the correct method
		params := ETLParams{
			// FIXME. This doesn't work.
			// We need to run this... every day on *yesterday*.
			Date: dateVal,
			Name: step,
		}
		// If anything fails, we exit the sequence
		// Pass nil so we don't write to the context in the loop
		c.Set("isSequence", true)
		err = put(c, params)
		if err != nil {
			errStep = step
			break
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": fmt.Sprintf("error in sequence execution: %s", errStep),
			"detail":  err.Error(),
			"date":    date,
			"name":    name,
		})
		return
	}
}

func seq_delete(c *gin.Context, seqPostbody SequencePostBody) {

}
