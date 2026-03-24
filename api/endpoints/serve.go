package serve

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	status "github.com/appleboy/gin-status-api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/csp"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

func Serve() {
	buffer_length := viper.GetInt("buffer.length")
	buffer_flushwaitsec := viper.GetInt("buffer.flushwaitsec")
	log.Printf("buffer length: %d flushwaitsec: %d\n", buffer_length, buffer_flushwaitsec)

	// Build the process network for buffering and
	// saving events that come in via the API
	ch_eb := make(chan types.EventBuffers)
	ch_evt := make(chan *types.Event)
	//ch_funs := make(chan func())

	// FIXME: add these constants to the configuration
	go csp.Enqueue(ch_evt, ch_eb, buffer_length, buffer_flushwaitsec)
	go csp.FlushBuffers(ch_eb)
	// This updates *yesterdays* database on minute one of the day

	engine := PourGin(ch_evt)

	host := viper.GetString("serve.host")
	port := viper.GetString("serve.port")
	_ = engine.Run(fmt.Sprintf("%s:%s", host, port))
}

func PourGin(ch_evt_out chan<- *types.Event) *gin.Engine {
	router := gin.Default()

	// We may want production mode.
	// This is configured via the envrionment
	if viper.GetString("ginmode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// See https://gin-gonic.com/en/docs/deployment/
	router.SetTrustedProxies(strings.Split(viper.GetString("proxies"), ","))

	log.Println("Setting default CORS handling")
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"PUT", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Api-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiV1 := router.Group("/v1")

	addMetadataEndpoints(apiV1)
	addLoggingEndpoints(apiV1, ch_evt_out)
	//addTestingEndpoints(apiV1, ch_evt_out)
	//addEtlEndpoints(apiV1)
	//addSequenceEndpoints(apiV1)

	return router
}

func addMetadataEndpoints(rG *gin.RouterGroup) {
	// The healthcheck is a public endpoint
	rG.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// The status endpoint provides server data, and needs a valid key
	// https://github.com/appleboy/gin-status-api
	auth_adminV1 := rG.Group("/")
	permissions := []types.PermissionType{
		constants.Admin,
		constants.Log,
		constants.ReadOnly,
		constants.Test,
	}
	auth_adminV1.Use(AuthMiddleWare(permissions))
	rG.GET("/status", status.GinHandler)

}
