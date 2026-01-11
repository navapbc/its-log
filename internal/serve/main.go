package serve

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/config"
	"github.com/jadudm/its-log/internal/csp"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/spf13/viper"
)

func PourGin(apiKeys config.ApiKeys, ch_evt_out chan<- *itslog.Event) *gin.Engine {
	// We may want production mode.
	// This is configured via the envrionment
	if viper.GetString("ginmode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	// See https://gin-gonic.com/en/docs/deployment/
	router.SetTrustedProxies(strings.Split(viper.GetString("proxies.trusted"), ","))

	// The entire API is currently v1
	apiV1 := router.Group("/v1")

	// Logging
	auth_logV1 := apiV1.Group("/")
	auth_logV1.Use(AuthMiddleWare(config.KEY_KIND_LOGGING, apiKeys))
	auth_logV1.PUT("se/:appID/:eventID", Event("se", ch_evt_out))
	auth_logV1.PUT("sev/:appID/:eventID/:value", Event("sev", ch_evt_out))
	auth_logV1.PUT("cse/:cluster/:appID/:eventID", Event("cse", ch_evt_out))
	auth_logV1.PUT("csev/:cluster/:appID/:eventID/:value", Event("csev", ch_evt_out))

	// ETL
	auth_adminV1 := apiV1.Group("/")
	auth_adminV1.Use(AuthMiddleWare(config.KEY_KIND_ADMIN, apiKeys))
	// Insert a new ETL step
	auth_adminV1.POST("etl/:date/:name", ETL)
	// Run an ETL step
	auth_adminV1.PUT("etl/:date/:name", ETL)
	// Retrieve the contents of a step, including the last run and run status
	auth_adminV1.GET("etl/:date/:name", ETL)
	// Remove a step
	auth_adminV1.DELETE("etl/:date/:name", ETL)

	// Querying the data
	auth_readV1 := apiV1.Group("/")
	auth_readV1.Use(AuthMiddleWare(config.KEY_KIND_READONLY, apiKeys))
	auth_readV1.GET("select/:date/:operation", Read)

	return router
}

func Serve(storage itslog.ItsLog, apiKeys config.ApiKeys) {
	buffer_length := viper.GetInt("buffer.length")
	buffer_flushwaitsec := viper.GetInt("buffer.flushwaitsec")
	log.Printf("buffer length: %d flushwaitsec: %d\n", buffer_length, buffer_flushwaitsec)
	// Build the process network for buffering and
	// saving events that come in via the API
	ch_eb := make(chan csp.EventBuffers)
	ch_evt := make(chan *itslog.Event)

	// FIXME: add these constants to the configuration
	go csp.Enqueue(ch_evt, ch_eb, buffer_length, buffer_flushwaitsec)
	go csp.FlushBuffers(ch_eb, storage)
	// This updates *yesterdays* database on minute one of the day

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
