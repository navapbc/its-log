package endpoints

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
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

	s := types.NewStorage(appId)
	err := s.SetDate(body.Date)
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

	err = s.Init()
	if err != nil {
		log.Println("storage init error: " + err.Error())
		panic(err)
	}
	// We cache whether this is loaded, so it is safe/fast to check every time
	// we try and load another ETL into the table.
	base.LoadDefaultEtlFiles(s)

	// Load an ItslogSummary structure first.
	// This has a hashing method that we can use to correctly hash
	// the summary on insert. This saves calling `hash-summaries`
	// after-the-fact.
	ils := models.ItslogSummary{
		LastRun:   time.Now(),
		KeyID:     keyId,
		Date:      body.Date,
		Operation: body.Operation,
		Tags:      body.Tags,
		Value:     body.Value,
		Count:     body.Count,
	}
	ils.HashItslogSummary()

	if err := s.Queries.InsertFullSummary(context.Background(), models.InsertFullSummaryParams{
		LastRun:   ils.LastRun,
		KeyID:     ils.KeyID,
		Date:      ils.Date,
		Operation: ils.Operation,
		Tags:      ils.Tags,
		Value:     ils.Value,
		Count:     ils.Count,
		Hash:      ils.Hash,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"method":  c.Request.Method,
			"message": "could not save summary",
			"date":    body.Date,
			"name":    body.Operation,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"method": c.Request.Method,
		"date":   body.Date,
		"name":   body.Operation,
	})
}
