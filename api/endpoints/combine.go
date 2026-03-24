package serve

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"io/fs"
// 	"log"
// 	"net/http"
// 	"os"
// 	"path/filepath"

// 	"github.com/gin-gonic/gin"
// 	"github.com/navapbc/its-log/internal/sqlite"
// 	"github.com/navapbc/its-log/internal/sqlite/models"
// 	"github.com/spf13/viper"
// )

// // https://stackoverflow.com/a/10510783
// func exists(path string) (bool, error) {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return true, nil
// 	}
// 	if errors.Is(err, fs.ErrNotExist) {
// 		return false, nil
// 	}
// 	return false, err
// }

// func Combine(c *gin.Context) {
// 	source := c.Param("source")
// 	root_path := viper.GetString("storage.path")

// 	// Quietly don't do anything if the source DB does not exist?
// 	exist, err := exists(filepath.Join(root_path, source+".sqlite"))
// 	if err != nil || !exist {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"message": fmt.Sprintf("%s does not exist", source+".sqlite"),
// 		})
// 		return
// 	}

// 	sourceStorage := &sqlite.SqliteStorage{
// 		Path:     viper.GetString("storage.path"),
// 		Filename: source + ".sqlite",
// 		KeyId:    GetKeyId(c),
// 	}
// 	sourceStorage.Init()
// 	defer sourceStorage.Close()

// 	destination := c.Param("destination")
// 	destinationStorage := &sqlite.SqliteStorage{
// 		Path:     viper.GetString("storage.path"),
// 		Filename: destination + ".sqlite",
// 		KeyId:    GetKeyId(c),
// 	}
// 	destinationStorage.Init()
// 	defer destinationStorage.Close()

// 	// FIXME: This is not yet designed to work on arbitrary tables
// 	// table := c.Param("table")

// 	rows, err := sourceStorage.GetQueries().GetAllSummaries(context.Background())
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "error",
// 			"message": "could not retrieve summaries",
// 		})
// 		return
// 	}

// 	for _, row := range rows {
// 		dctx := context.Background()
// 		err := destinationStorage.GetQueries().InsertSummary(dctx, models.InsertSummaryParams{
// 			KeyID:      row.KeyID,
// 			Date:       row.Date,
// 			Operation:  row.Operation,
// 			SourceName: row.SourceName,
// 			EventName:  row.EventName,
// 			Value:      row.Value,
// 		})
// 		if err != nil {
// 			log.Println("InsertSummary error: " + fmt.Sprintf("%s", err))
// 		}
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"status":  "ok",
// 		"message": "rows copied: " + fmt.Sprintf("%d", len(rows)),
// 	})
// }
