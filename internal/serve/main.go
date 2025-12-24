package serve

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/config"
	"github.com/jadudm/its-log/internal/csp"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/spf13/viper"
)

func PourGin(apiKeys config.ApiKeys, ch_evt_out chan<- *itslog.Event) *gin.Engine {
	// We may want production mode.
	// This is configured via the envrionment
	if viper.GetString("gin_mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	apiV1 := router.Group("/v1")
	authV1 := apiV1.Group("/")
	authV1.Use(AuthMiddleWare(apiKeys))

	authV1.PUT("event/:appID/:eventID", Event(ch_evt_out))

	return router
}

func Serve(storage itslog.ItsLog, apiKeys config.ApiKeys) {
	event_buffer_length := viper.GetInt("app.event_buffer_length")
	event_buffer_flush_seconds := viper.GetInt("app.event_buffer_flush_seconds")

	// Build the process network for buffering and
	// saving events that come in via the API
	ch_eb := make(chan csp.EventBuffers)
	ch_evt := make(chan *itslog.Event)

	// FIXME: add these constants to the configuration
	go csp.Enqueue(ch_evt, ch_eb, event_buffer_length, event_buffer_flush_seconds)
	go csp.FlushBuffers(ch_eb, storage)
	// This updates *yesterdays* database on minute one of the day
	go csp.Summarize()

	engine := PourGin(apiKeys, ch_evt)
	host := viper.GetString("serve.host")
	port := viper.GetString("serve.port")
	cert := viper.GetString("serve.cert")
	key := viper.GetString("serve.key")
	if cert != "mock" && key != "mock" {
		_ = engine.RunTLS(fmt.Sprintf("%s:%s", host, port), cert, key)
	}
	panic("failed to find cert/key. leaving in a panic.")
}
