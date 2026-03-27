package endpoints

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/schema/models"
	"github.com/navapbc/its-log/internal/types"
)

type SummaryReadBody struct {
	Date      string `json:"date" binding:"required"`
	Tags      string `json:"tags" binding:"required"`
	Operation string `json:"operation" binding:"required"`
}

func SummaryRead(c *gin.Context) {
	var body SummaryReadBody
	// Call ShouldBindJSON to bind the incoming JSON to the newItem struct
	if err := c.ShouldBindJSON(&body); err != nil {
		// If an error occurs (e.g., invalid JSON, missing required fields),
		// return a 400 Bad Request error
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("should bind: %s", err.Error())})
		return
	}

	appId := base.GetOrPanic(c, "AppId")

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

	params := models.ReadSummaryParams{
		Tags:      body.Tags,
		Operation: body.Operation,
	}

	row, err := s.Queries.ReadSummary(context.Background(), params)
	if err != nil {
		log.Println("err: " + err.Error())
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    constants.OK,
		"operation": body.Operation,
		"tags":      body.Tags,
		"count":     row.Count,
	})

}
