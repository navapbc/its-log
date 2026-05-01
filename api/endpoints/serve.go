package endpoints

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"

	status "github.com/appleboy/gin-status-api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/csp"
	"github.com/navapbc/its-log/internal/types"
	"github.com/spf13/viper"
)

func Serve(params types.ServeParams) {
	buffer_length := viper.GetInt("buffer.length")
	buffer_flushwaitsec := viper.GetInt("buffer.flushwaitsec")
	log.Printf("buffer length: %d flushwaitsec: %d\n", buffer_length, buffer_flushwaitsec)

	// Build the process network for buffering and
	// saving events that come in via the API
	ch_eb := make(chan types.EventBuffer)
	ch_evt := make(chan *types.Event)
	//ch_funs := make(chan func())

	// FIXME: add these constants to the configuration
	go csp.Enqueue(ch_evt, ch_eb, buffer_length, buffer_flushwaitsec)
	go csp.FlushBuffers(ch_eb)
	// This updates *yesterdays* database on minute one of the day

	engine := PourGin(ch_evt)
	if params.Mode == "debug" {
		pprof.Register(engine, "dev/pprof")
	}

	host := viper.GetString("serve.host")
	port := viper.GetString("serve.port")
	_ = engine.Run(fmt.Sprintf("%s:%s", host, port))
}

func PourGin(ch_evt_out chan<- *types.Event) *gin.Engine {
	router := gin.Default()

	// Adding this will override the default logging in Gin.
	// May be useful to flip entirely to zap.
	// Also https://stackoverflow.com/questions/48780070/disable-request-logging-for-a-particular-go-gin-route
	// gin.DefaultWriter = io.Discard

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

	// https://github.com/gin-contrib/zap?tab=readme-ov-file#example
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))

	apiV1 := router.Group("/v1")

	addLoggingEndpoints(apiV1, ch_evt_out)
	addMetadataEndpoints(apiV1)
	//addTestingEndpoints(apiV1, ch_evt_out)
	addEtlEndpoints(apiV1)
	addSequenceEndpoints(apiV1)
	addSummaryEndpoints(apiV1)

	return router
}

// LOGGING ENDPOINTS
func addLoggingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *types.Event) {
	// Logging
	auth_logV1 := rG.Group("/")
	permissions := []types.PermissionType{constants.Log, constants.Test}
	auth_logV1.Use(AuthMiddleWare(permissions))
	auth_logV1.POST(constants.LOG_CREATE, LogCreate(ch_evt_out, constants.Log))
	auth_logV1.POST(constants.LOG_CREATE_DATE, LogCreate(ch_evt_out, constants.Test))
}

// ETL ENDPOINTS
func addEtlEndpoints(rG *gin.RouterGroup) {
	auth_adminV1 := rG.Group("/")
	permissions := []types.PermissionType{constants.Admin, constants.Test}
	auth_adminV1.Use(AuthMiddleWare(permissions))

	// Insert a new ETL step
	auth_adminV1.POST(constants.ETL_CREATE, CreateEtl)
	// Run an ETL step
	auth_adminV1.POST(constants.ETL_RUN, RunEtl)
}

// METADATA ENDPOINTS
func addMetadataEndpoints(rG *gin.RouterGroup) {
	// The healthcheck is a public endpoint
	rG.GET(constants.METADATA_HEALTH, func(c *gin.Context) {
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
	auth_adminV1.GET(constants.METADATA_STATUS, status.GinHandler)

}

// SEQUENCE ENDPOINTS
func addSequenceEndpoints(rG *gin.RouterGroup) {
	auth_adminV1 := rG.Group("/")
	permissions := []types.PermissionType{constants.Admin, constants.Test}
	auth_adminV1.Use(AuthMiddleWare(permissions))

	// Insert a sequence
	// It's actually just inserting an entry into the
	// ETL table with the correct values. We give it a 'sequence'
	// endpoint, but it is interchangeable.
	auth_adminV1.POST(constants.SEQUENCE_CREATE, CreateEtl)
	// Run a sequence
	auth_adminV1.POST(constants.SEQUENCE_RUN, RunSequence)
}

// SUMMARY ENDPOINTS
func addSummaryEndpoints(rG *gin.RouterGroup) {
	auth_adminV1ReadOnly := rG.Group("/")
	auth_adminV1ReadWrite := rG.Group("/")

	permsRO := []types.PermissionType{constants.Admin, constants.ReadOnly, constants.Test}
	permsRW := []types.PermissionType{constants.Admin, constants.Test}
	auth_adminV1ReadOnly.Use(AuthMiddleWare(permsRO))
	auth_adminV1ReadWrite.Use(AuthMiddleWare(permsRW))

	auth_adminV1ReadOnly.POST(constants.SUMMARY_READ, SummaryRead)
	auth_adminV1ReadWrite.POST(constants.SUMMARY_CREATE, SummaryCreate)

}
