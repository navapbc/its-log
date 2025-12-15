package serve

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/itslog"
)

func Event(storage itslog.ItsLog) func(c *gin.Context) {
	return func(c *gin.Context) {
		appID := c.Param("appID")
		evtID := c.Param("eventID")
		storage.Event(&itslog.Event{
			Source: appID,
			Event:  evtID,
		})
		// Everything worked.
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})

	}

}
