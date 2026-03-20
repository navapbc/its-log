package serve

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jadudm/its-log/internal/base"
)

// @BasePath /v1
// Thursday, 2PM?

func addLoggingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *base.Event) {
	// Logging
	auth_logV1 := rG.Group("/")
	permissions := []base.PermissionType{base.Log, base.Test}
	auth_logV1.Use(AuthMiddleWare(permissions))
	auth_logV1.POST("/log", Event(ch_evt_out, base.Log))
}

func setTimestamp(c *gin.Context, evt *base.Event, permission base.PermissionType) {
	switch permission {
	case base.Test:
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

	case base.Log:
		evt.Timestamp = time.Now()

	default:
		// If we manage to get here, make sure we set the timestamp to
		// the current time as a reasonable default behavior.
		evt.Timestamp = time.Now()
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
func Event(ch_evt_out chan<- *base.Event, permission base.PermissionType) func(c *gin.Context) {
	return func(c *gin.Context) {
		// https://pkg.go.dev/github.com/go-playground/validator/v10
		var evt base.Event
		// Call ShouldBindJSON to parse the request body into the struct
		if err := c.ShouldBindJSON(&evt); err != nil {
			// FIXME: follow standard response protocol
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		evt.AppId = base.GetOrPanic(c, base.ITSLOG_APPID)
		evt.KeyId = base.GetOrPanic(c, base.ITSLOG_KEYID)

		// If it is a test event, we mangle a date parameter.
		// If it is not a test event, we use Now().
		setTimestamp(c, &evt, permission)

		// The tags field is currently a JSON array. It needs to become
		// a sorted, dot-separated string.
		sort.Strings(evt.Tags)
		evt.TagString = strings.Join(evt.Tags, ".")

		// This uses the struct validator library, which provides
		// a rich notion of contracts over the fields in the struct
		err := validate.Struct(evt)
		if err != nil {
			messages := make([]string, 0)
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				for _, fieldError := range validationErrors {
					messages = append(messages,
						fmt.Sprintf("Field: %s, Tag: %s",
							fieldError.Field(), fieldError.Tag()))

				}
			}

			c.JSON(http.StatusInternalServerError, base.Error{
				Status:    base.ERROR,
				Error:     fmt.Sprintf("%s", messages),
				ErrorType: "field validation errors",
			})

			return
		}

		// Send the event to the Enqueue-er
		ch_evt_out <- &evt

		// Everything worked.
		c.JSON(http.StatusOK, base.Success{Status: base.OK})
	}
}
