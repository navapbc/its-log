package serve

import (
	"fmt"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/itslog"
)

func addTestingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *itslog.Event) {
	// Test data generation
	auth_testV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Test}
	auth_testV1.Use(AuthMiddleWare(permissions))
	auth_testV1.PUT("dse/:date/:source/:event", TestEvent("dse", ch_evt_out))
	auth_testV1.PUT("dsev/:date/:source/:event/:value", TestEvent("dsev", ch_evt_out))
	auth_testV1.PUT("dcse/:date/:cluster/:source/:event", TestEvent("dse", ch_evt_out))
	auth_testV1.PUT("dcsev/:date/:cluster/:appID/:eventID/:value", TestEvent("dse", ch_evt_out))
}

func getTestEventType(root string) itslog.EventType {
	mapping := map[string]itslog.EventType{
		"dse":   itslog.DSE,
		"dsev":  itslog.DSEV,
		"dcse":  itslog.DCSE,
		"dcsev": itslog.DCSEV,
	}
	// As written, this will panic hard if we get a value we're not expecting.
	// However, these should be our own endpoints/the code in `addTestingEndpoints`
	// so the chances of a panic are low.
	return mapping[root]
}

func TestEvent(root string, ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
	return func(c *gin.Context) {
		// https://pkg.go.dev/github.com/go-playground/validator/v10

		eventType := getEventType(root)
		cluster := ""
		source := c.Param("source")
		event := c.Param("event")
		value := ""

		appId := itslog.GetOrPanic(c, "AppId")
		keyId := itslog.GetOrPanic(c, "KeyId")

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

		// Only some event types will require the cluster
		if slices.Contains([]itslog.EventType{itslog.DCSE, itslog.DCSEV}, eventType) {
			cluster = c.Param("cluster")
		}

		// And only some have a value
		if slices.Contains([]itslog.EventType{itslog.DSEV, itslog.DCSEV}, eventType) {
			value = c.Param("value")
		}

		payload := &itslog.Event{
			Timestamp: timestamp,
			KeyId:     keyId,
			AppId:     appId,
			Cluster:   cluster,
			Source:    source,
			Event:     event,
			Value:     value,
		}

		err = validate.Struct(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "could not validate event",
			})
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

// // Event handling has become overloaded, and should be simplified.
// func Event(root string, ch_evt_out chan<- *itslog.Event) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		var timestamp time.Time

// 		evtID := c.Param("eventId")
// 		sourceId := c.Param("sourceId")

// 		switch  {
// 			case

// 		}

// 		cluster := ""
// 		value := ""

// 		appId := itslog.GetOrPanic("AppId")

// 		if strings.Contains(root, "d") {
// 			var err error
// 			date := c.Param("date")

// 			timestamp, err = time.Parse("2006-01-02", date)
// 			min := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC).Unix()
// 			max := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 23, 59, 59, 0, time.UTC).Unix()
// 			delta := max - min
// 			sec := rand.Int63n(delta) + min
// 			timestamp = time.Unix(sec, 0)

// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"status":  "error",
// 					"message": fmt.Sprintf("date is malformed: %s", date),
// 				})
// 				return
// 			}
// 		} else {
// 			timestamp = time.Now()
// 		}
// 		if strings.Contains(root, "c") {
// 			cluster = c.Param("cluster")
// 		}
// 		if strings.Contains(root, "v") {
// 			value = c.Param("value")
// 		}

// 		// Send the event to the Enqueue-er
// 		ch_evt_out <- &itslog.Event{
// 			Timestamp: timestamp,
// 			AppId:     appId,
// 			Cluster:   cluster,
// 			Source:    appID,
// 			Event:     evtID,
// 			Value:     value,
// 		}

// 		// Everything worked.
// 		c.JSON(http.StatusOK, gin.H{
// 			"status": "ok",
// 		})
// 	}
// }
