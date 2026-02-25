package serve

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/itslog"
)

func addTestingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *itslog.Event) {
	// Test data generation
	auth_testV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Test}
	auth_testV1.Use(AuthMiddleWare(permissions))
	auth_testV1.POST("log/:date", TestEvent(ch_evt_out))
}

func TestEvent(ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return func(c *gin.Context) {
		// https://pkg.go.dev/github.com/go-playground/validator/v10
		var evt *itslog.Event
		// Call ShouldBindJSON to parse the request body into the struct
		if err := c.ShouldBindJSON(evt); err != nil {
			// FIXME: follow standard response protocol
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		evt.AppId = itslog.GetOrPanic(c, "AppId")
		evt.KeyId = itslog.GetOrPanic(c, "KeyId")

		// The date in the testing cases comes from the URL.
		// Hence, we might have parsing errors on what is passed in.
		var err error
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

		err = validate.Struct(evt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "could not validate event",
			})
		}

		// Send the event to the Enqueue-er
		ch_evt_out <- evt
		// Everything worked.
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "its logged",
		})
	}
}
