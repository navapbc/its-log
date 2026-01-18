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

// @BasePath /v1

func addLoggingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *itslog.Event) {
	// Logging
	auth_logV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Log, itslog.Test}
	auth_logV1.Use(AuthMiddleWare(permissions))
	auth_logV1.PUT("se/:source/:event", EventSE(ch_evt_out))
	auth_logV1.PUT("sev/:source/:event/:value", EventSEV(ch_evt_out))
	auth_logV1.PUT("cse/:cluster/:source/:event", EventCSE(ch_evt_out))
	auth_logV1.PUT("csev/:cluster/:source/:event/:value", EventCSEV(ch_evt_out))
}

// EventSE godoc
// @Accept json
// @Description logs an event with a source
// @Failure 500	{object} itslog.Error
// @Param event path string true "an event category"
// @Param source path string true "a source category"
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /se/{source}/{event} [put]
// @Schemes
// @Success 200 {object} itslog.Success
// @Summary log a source/event
// @Tags events
func EventSE(ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return Event(itslog.SE, ch_evt_out)
}

// EventSEV godoc
// @Accept json
// @Description log a source and event with a unique value
// @Failure 500	{object} itslog.Error
// @Param event path string true "an event category"
// @Param source path string true "a source category"
// @Param value path string true "a unique value associated with this event"
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /sev/{source}/{event}/{value} [put]
// @Schemes
// @Success 200 {object} itslog.Success
// @Summary log a source and event with a unique value
// @Tags events
func EventSEV(ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return Event(itslog.SEV, ch_evt_out)
}

// EventCSE godoc
// @Accept json
// @Description log a source and event as part of a cluster
// @Failure 500	{object} itslog.Error
// @Param cluster path string true "a UUID representing this cluster"
// @Param event path string true "an event category"
// @Param source path string true "a source category"
// @Param x-api-key header string true "API key, 32 bytes or more, issued"
// @Produce json
// @Router /cse/{cluster}/{source}/{event} [put]
// @Schemes
// @Success 200 {object} itslog.Success
// @Summary log a source and event as part of a cluster
// @Tags events
func EventCSE(ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return Event(itslog.CSE, ch_evt_out)
}

// param name,param type,data type,is mandatory?,comment attribute(optional)

// EventCSEV godoc
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
func EventCSEV(ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return Event(itslog.CSEV, ch_evt_out)
}

func Event(eventType itslog.EventType, ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return func(c *gin.Context) {
		// https://pkg.go.dev/github.com/go-playground/validator/v10
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

			c.JSON(http.StatusInternalServerError, itslog.Error{
				Status:    itslog.ERROR,
				Error:     fmt.Sprintf("%s", messages),
				ErrorType: "field validation errors",
			})

			return
		}

		// Send the event to the Enqueue-er
		ch_evt_out <- payload

		// Everything worked.
		c.JSON(http.StatusOK, itslog.Success{Status: itslog.OK})
	}
}
