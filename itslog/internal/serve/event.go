package serve

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/itslog"
)

// Event handling has become overloaded, and should be simplified.
func Event(root string, ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return func(c *gin.Context) {
		var timestamp time.Time

		appID := c.Param("appID")
		evtID := c.Param("eventID")
		cluster := ""
		value := ""

		if strings.Contains(root, "d") {
			var err error
			date := c.Param("date")

			timestamp, err = time.Parse("2006-01-02", date)
			min := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC).Unix()
			max := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 23, 59, 59, 0, time.UTC).Unix()
			delta := max - min
			sec := rand.Int63n(delta) + min
			timestamp = time.Unix(sec, 0)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("date is malformed: %s", date),
				})
				return
			}
		} else {
			timestamp = time.Now()
		}
		if strings.Contains(root, "c") {
			cluster = c.Param("cluster")
		}
		if strings.Contains(root, "v") {
			value = c.Param("value")
		}

		// This should have been handled by the auth middleware
		key_id, err := GetKeyId(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "no key id in the Gin context",
			})
			return
		}

		// Send the event to the Enqueue-er
		ch_evt_out <- &itslog.Event{
			Timestamp: timestamp,
			KeyId:     key_id,
			Cluster:   cluster,
			Source:    appID,
			Event:     evtID,
			Value:     value,
		}

		// Everything worked.
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}
