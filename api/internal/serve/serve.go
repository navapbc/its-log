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
	"github.com/go-playground/validator/v10"
	"github.com/jadudm/its-log/docs"
	"github.com/jadudm/its-log/internal/csp"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var validate *validator.Validate

func Serve() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	buffer_length := viper.GetInt("buffer.length")
	buffer_flushwaitsec := viper.GetInt("buffer.flushwaitsec")
	log.Printf("buffer length: %d flushwaitsec: %d\n", buffer_length, buffer_flushwaitsec)

	// Build the process network for buffering and
	// saving events that come in via the API
	ch_eb := make(chan csp.EventBuffers)
	ch_evt := make(chan *itslog.Event)

	// FIXME: add these constants to the configuration
	go csp.Enqueue(ch_evt, ch_eb, buffer_length, buffer_flushwaitsec)
	go csp.FlushBuffers(ch_eb)
	// This updates *yesterdays* database on minute one of the day

	engine := PourGin(ch_evt)

	host := viper.GetString("serve.host")
	port := viper.GetString("serve.port")
	_ = engine.Run(fmt.Sprintf("%s:%s", host, port))
}

func PourGin(ch_evt_out chan<- *itslog.Event) *gin.Engine {
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

	addSwaggerEndpoints(apiV1)
	addMetadataEndpoints(apiV1)
	addLoggingEndpoints(apiV1, ch_evt_out)
	addTestingEndpoints(apiV1, ch_evt_out)
	addEtlEndpoints(apiV1)
	addQueryEndpoints(apiV1)

	return router
}

func addSwaggerEndpoints(rG *gin.RouterGroup) {
	docs.SwaggerInfo.BasePath = rG.BasePath()
	rG.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

}

func addMetadataEndpoints(rG *gin.RouterGroup) {
	// The healthcheck is a public endpoint
	rG.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// The status endpoint provides server data, and needs a valid key
	// https://github.com/appleboy/gin-status-api
	auth_adminV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Admin, itslog.Log, itslog.ReadOnly, itslog.Test}
	auth_adminV1.Use(AuthMiddleWare(permissions))
	rG.GET("/status", status.GinHandler)

}

func addQueryEndpoints(rG *gin.RouterGroup) {
	// Fixme: rethink querying
	// Querying the data
	// auth_readV1 := rG.Group("/")
	// permissions := []itslog.PermissionType{}
	// auth_readV1.Use(AuthMiddleWare(permissions, apiKeys))
	// auth_readV1.GET("select/:date/:operation/:source_name", Read)

}
