package serve

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/b2"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/types"
)

// @BasePath /v1
// Thursday, 2PM?

func addLoggingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *types.Event) {
	// Logging
	auth_logV1 := rG.Group("/")
	permissions := []types.PermissionType{constants.Log, constants.Test}
	auth_logV1.Use(AuthMiddleWare(permissions))
	auth_logV1.POST("/log", Event(ch_evt_out, constants.Log))
}

func setTimestamp(c *gin.Context, evt *types.Event, permission types.PermissionType) {
	switch permission {
	case constants.Test:
		// The date in the testing cases comes from the URL.
		// Hence, we might have parsing errors on what is passed in.
		date := c.Param("date")
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

// Event godoc
// @Accept json
// @Description log a source and event with a unique value as part of a cluster
// @Failure 500	{object} base.Error
// @Param cluster path string true "a UUID representing this cluster"
// @Param event path string true "an event category"
// @Param source path string true "a source category"
// @Param value path string true "a unique value associated with this event"
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /csev/{cluster}/{source}/{event}/{value} [put]
// @Schemes
// @Success 200 {object} base.Success
// @Summary log a source and event with a unique value as part of a cluster
// @Tags events
func Event(ch_evt_out chan<- *types.Event, permission types.PermissionType) func(c *gin.Context) {
	return func(c *gin.Context) {
		// https://pkg.go.dev/github.com/go-playground/validator/v10
		var evt types.Event
		// Call ShouldBindJSON to parse the request body into the struct
		if err := c.ShouldBindJSON(&evt); err != nil {
			// FIXME: follow standard response protocol
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		evt.AppId = b2.GetOrPanic(c, constants.ITSLOG_APPID)
		evt.KeyId = b2.GetOrPanic(c, constants.ITSLOG_KEYID)

		// If it is a test event, we mangle a date parameter.
		// If it is not a test event, we use Now().
		setTimestamp(c, &evt, permission)

		// The tags field is currently a JSON array. It needs to become
		// a sorted, dot-separated string.
		sort.Strings(evt.Tags)
		evt.TagString = strings.Join(evt.Tags, ".")

		// Send the event to the Enqueue-er
		ch_evt_out <- &evt

		// Everything worked.
		c.JSON(http.StatusOK, types.Success{Status: constants.OK})
	}
}
