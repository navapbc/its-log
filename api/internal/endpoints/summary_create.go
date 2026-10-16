package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/storage"
	"github.com/navapbc/its-log/internal/types"
	"go.uber.org/zap"
)

func SummaryCreate(c *gin.Context) {

	var body types.SummaryPostBody
	// Call ShouldBindJSON to bind the incoming JSON to the newItem struct
	if err := c.ShouldBindJSON(&body); err != nil {
		// If an error occurs (e.g., invalid JSON, missing required fields),
		// return a 400 Bad Request error
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("should bind: %s", err.Error())})
		return
	}

	appId := base.GetOrPanic(c, "AppId")
	keyId := base.GetOrPanic(c, "KeyId")

	s, err := storage.NewStorage(appId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not initialize storage backend",
			"date":    body.Date,
			"name":    appId,
		})
		return
	}

	err = s.SetDateYMD(body.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not parse date; must be YYYY-MM-DD",
			"date":    body.Date,
			"name":    body.Operation,
		})
		return
	}

	err = s.InitDB()
	if err != nil {
		log.Println("storage init error: " + err.Error())
		panic(err)
	}

	// The tags field comes in as a JSON array. It needs to become
	// a sorted, dot-separated string.
	sort.Strings(body.Tags)
	theTags := strings.Join(body.Tags, ".")

	_, err = s.Exec(storage.Endpoint, "insert_full_summary", storage.EmbedParamMap{
		"key_id":    keyId,
		"date":      body.Date,
		"operation": body.Operation,
		"tags":      theTags,
		"value":     body.Value,
		"count":     body.Count,
	})

	if err != nil {
		zap.L().Error("summary save error", zap.Error(err))

		c.JSON(http.StatusNotImplemented, gin.H{
			"status":    "error",
			"method":    c.Request.Method,
			"message":   "could not save ETL step",
			"date":      body.Date,
			"operation": body.Operation,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":    "ok",
		"method":    c.Request.Method,
		"date":      body.Date,
		"operation": body.Operation,
		"tags":      theTags,
		"value":     body.Value,
		"count":     body.Count,
	})

}
