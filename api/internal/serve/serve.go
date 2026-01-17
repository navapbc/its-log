package serve

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jadudm/its-log/internal/csp"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/spf13/viper"
)

var validate *validator.Validate

func configureEngine() *gin.Engine {
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

	return engine
}

func Serve() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	engine := configureEngine()

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
	// router.SetTrustedProxies(strings.Split(viper.GetString("proxies.trusted"), ","))
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

	addHealthCheck(apiV1)
	addLoggingEndpoints(apiV1, ch_evt_out)
	addTestingEndpoints(apiV1, ch_evt_out)
	//addEtlEndpoints(apiV1, ch_evt_out)
	addQueryEndpoints(apiV1, ch_evt_out)

	return router
}

func addHealthCheck(rG *gin.RouterGroup) {
	rG.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func addQueryEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *itslog.Event) {
	// Fixme: rethink querying
	// Querying the data
	// auth_readV1 := rG.Group("/")
	// permissions := []itslog.PermissionType{}
	// auth_readV1.Use(AuthMiddleWare(permissions, apiKeys))
	// auth_readV1.GET("select/:date/:operation/:source_name", Read)

}
