package serve

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/itslog"
)

func Event(root string, ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return func(c *gin.Context) {
		appID := c.Param("appID")
		evtID := c.Param("eventID")
		cluster := ""
		value := ""

		if strings.Contains(root, "c") {
			cluster = c.Param("cluster")
		}
		if strings.Contains(root, "v") {
			value = c.Param("value")
		}

		// Send the event to the Enqueue-er
		ch_evt_out <- &itslog.Event{
			Timestamp: time.Now(),
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
