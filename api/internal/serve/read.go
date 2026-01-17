package serve

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jadudm/its-log/internal/sqlite"
// 	"github.com/jadudm/its-log/internal/sqlite/models"
// 	"github.com/spf13/viper"
// )

// func Read(c *gin.Context) {
// 	filename := c.Param("date")
// 	source_name := c.Param("source_name")
// 	operation := c.Param("operation")
// 	keyId := GetKeyId(c)

// 	s := sqlite.SqliteStorage{
// 		Path:     viper.GetString("storage.path"),
// 		Filename: fmt.Sprintf("%s.sqlite", filename),
// 		KeyId:    keyId,
// 	}

// 	s.Init()

// 	rows, err := s.GetQueries().ReadSummary(context.Background(), models.ReadSummaryParams{
// 		KeyID:      keyId,
// 		Operation:  operation,
// 		SourceName: sql.NullString{Valid: true, String: source_name},
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  "error",
// 			"message": "could not read summaries: " + err.Error(),
// 		})
// 		return
// 	}

// 	log.Printf("Retrieved %d rows\n", len(rows))

// 	// Everything worked.
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":      "ok",
// 		"date":        filename,
// 		"operation":   operation,
// 		"source_name": source_name,
// 		"results":     rows,
// 	})
// }
