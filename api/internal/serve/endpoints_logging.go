package serve

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jadudm/its-log/internal/itslog"
)

func addLoggingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *itslog.Event) {
	// Logging
	auth_logV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Logging, itslog.Test}
	auth_logV1.Use(AuthMiddleWare(permissions))
	auth_logV1.PUT("se/:source/:event", Event("se", ch_evt_out))
	auth_logV1.PUT("sev/:source/:event/:value", Event("sev", ch_evt_out))
	auth_logV1.PUT("cse/:cluster/:source/:event", Event("cse", ch_evt_out))
	auth_logV1.PUT("csev/:cluster/:source/:event/:value", Event("csev", ch_evt_out))
}

func getEventType(root string) itslog.EventType {
	mapping := map[string]itslog.EventType{
		"se":   itslog.SE,
		"sev":  itslog.SEV,
		"cse":  itslog.CSE,
		"csev": itslog.CSEV,
	}
	// As written, this will panic hard if we get a value we're not expecting.
	// However, these should be our own endpoints/the code in `addLoggingEndpoints`
	// so the chances of a panic are low.
	return mapping[root]
}

func Event(root string, ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return func(c *gin.Context) {
		// https://pkg.go.dev/github.com/go-playground/validator/v10

		eventType := getEventType(root)
		timestamp := time.Now()
		cluster := ""
		source := c.Param("source")
		event := c.Param("event")
		value := ""

		appId := itslog.GetOrPanic(c, "AppId")
		keyId := itslog.GetOrPanic(c, "KeyId")

		// Only some event types will require the cluster
		if slices.Contains([]itslog.EventType{itslog.CSE, itslog.CSEV}, eventType) {
			cluster = c.Param("cluster")
		}

		// And only some have a value
		if slices.Contains([]itslog.EventType{itslog.SEV, itslog.CSEV}, eventType) {
			value = c.Param("value")
		}

		payload := &itslog.Event{
			Timestamp: timestamp,
			EventType: eventType,
			KeyId:     keyId,
			AppId:     appId,
			Cluster:   cluster,
			Source:    source,
			Event:     event,
			Value:     value,
		}

		// This uses the struct validator library, which provides
		// a rich notion of contracts over the fields in the struct
		err := validate.Struct(payload)
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

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("field validation errors: %s", messages),
			})
			return
		}

		// Send the event to the Enqueue-er
		ch_evt_out <- payload

		// Everything worked.
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "its logged",
		})
	}
}
