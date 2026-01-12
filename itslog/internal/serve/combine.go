package serve

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/sqlite"
	"github.com/jadudm/its-log/internal/sqlite/models"
	"github.com/spf13/viper"
)

func Combine(c *gin.Context) {
	source := c.Param("source")
	sourceStorage := &sqlite.SqliteStorage{
		Path: viper.GetString("storage.path"),
	}
	sourceStorage.InitByName(source)
	defer sourceStorage.Close()

	destination := c.Param("destination")
	destinationStorage := &sqlite.SqliteStorage{
		Path: viper.GetString("storage.path"),
	}
	destinationStorage.InitByName(destination)
	defer destinationStorage.Close()

	// FIXME: This is not yet designed to work on arbitrary tables
	// table := c.Param("table")

	rows, err := sourceStorage.GetQueries().GetAllSummaries(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "could not retrieve summaries",
		})
		return
	}

	// 	type ItslogSummary struct {
	// 	ID         int64
	// 	KeyID      string
	// 	Date       time.Time
	// 	Operation  string
	// 	SourceName sql.NullString
	// 	EventName  sql.NullString
	// 	Value      float64
	// }

	for _, row := range rows {
		dctx := context.Background()
		fmt.Printf("%v\n", row)
		err := destinationStorage.GetQueries().InsertSummary(dctx, models.InsertSummaryParams{
			KeyID:      row.KeyID,
			Date:       row.Date,
			Operation:  row.Operation,
			SourceName: row.SourceName,
			EventName:  row.EventName,
			Value:      row.Value,
		})
		if err != nil {
			log.Println("InsertSummary error: " + fmt.Sprintf("%s", err))
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "ok",
		"message": "rows copied: " + fmt.Sprintf("%d", len(rows)),
	})
}
