package endpoints

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/types"
)

func LogCreate(ch_evt_out chan<- *types.Event, permission types.PermissionType) func(c *gin.Context) {
	return func(c *gin.Context) {
		var evt types.Event
		if err := c.ShouldBindJSON(&evt); err != nil {
			// FIXME: follow standard response protocol
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		evt.AppId = base.GetOrPanic(c, constants.ITSLOG_APPID)
		evt.KeyId = base.GetOrPanic(c, constants.ITSLOG_KEYID)

		setTimestamp(c, &evt, permission)

		// The tags field comes in as a JSON array. It needs to become
		// a sorted, dot-separated string.
		sort.Strings(evt.Tags)
		evt.TagString = strings.Join(evt.Tags, ".")

		// Send the event to the Enqueue-er
		ch_evt_out <- &evt

		// Everything worked.
		c.JSON(http.StatusOK, types.Success{Status: constants.OK})
	}
}

func setTimestamp(c *gin.Context, evt *types.Event, permission types.PermissionType) {
	switch permission {
	case constants.Test:
		// If this is a test event, it will come in on the /create/date pathway, and
		// have a `date` field in the JSON, and the permission will be `Test`.
		date := evt.Date
		timestamp, err := time.Parse("2006-01-02", date)
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

		evt.Timestamp = timestamp

	case constants.Log:
		evt.Timestamp = time.Now().UTC()

	default:
		// If we manage to get here, make sure we set the timestamp to
		// the current time as a reasonable default behavior.
		evt.Timestamp = time.Now().UTC()
	}

}
