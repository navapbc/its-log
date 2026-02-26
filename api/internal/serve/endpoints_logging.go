package serve

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jadudm/its-log/internal/itslog"
)

// @BasePath /v1
// Thursday, 2PM?

func addLoggingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *itslog.Event) {
	// Logging
	auth_logV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Log, itslog.Test}
	auth_logV1.Use(AuthMiddleWare(permissions))
	auth_logV1.POST("/log", Event(ch_evt_out))
}

// Event godoc
// @Accept json
// @Description log a source and event with a unique value as part of a cluster
// @Failure 500	{object} itslog.Error
// @Param cluster path string true "a UUID representing this cluster"
// @Param event path string true "an event category"
// @Param source path string true "a source category"
// @Param value path string true "a unique value associated with this event"
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /csev/{cluster}/{source}/{event}/{value} [put]
// @Schemes
// @Success 200 {object} itslog.Success
// @Summary log a source and event with a unique value as part of a cluster
// @Tags events
func Event(ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return func(c *gin.Context) {
		// https://pkg.go.dev/github.com/go-playground/validator/v10
		var evt itslog.Event
		// Call ShouldBindJSON to parse the request body into the struct
		if err := c.ShouldBindJSON(&evt); err != nil {
			// FIXME: follow standard response protocol
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		evt.Timestamp = time.Now()

		evt.AppId = itslog.GetOrPanic(c, "AppId")
		evt.KeyId = itslog.GetOrPanic(c, "KeyId")

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

			c.JSON(http.StatusInternalServerError, itslog.Error{
				Status:    itslog.ERROR,
				Error:     fmt.Sprintf("%s", messages),
				ErrorType: "field validation errors",
			})

			return
		}

		// Send the event to the Enqueue-er
		ch_evt_out <- &evt

		// Everything worked.
		c.JSON(http.StatusOK, itslog.Success{Status: itslog.OK})
	}
}
