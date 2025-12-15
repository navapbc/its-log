package serve

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/config"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/spf13/viper"
)

func PourGin(s itslog.ItsLog, apiKeys config.ApiKeys) *gin.Engine {
	// We may want production mode.
	// This is configured via the envrionment
	if viper.GetString("gin_mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	apiV1 := router.Group("/v1")
	authV1 := apiV1.Group("/")
	authV1.Use(AuthMiddleWare(apiKeys))

	authV1.PUT("event/:appID/:eventID", Event(s))
	authV1.PUT("unique/:appID/:eventID", Event(s))

	return router
}

func Serve(storage itslog.ItsLog, apiKeys config.ApiKeys) {
	engine := PourGin(storage, apiKeys)
	host := viper.GetString("serve.host")
	port := viper.GetString("serve.port")
	_ = engine.Run(fmt.Sprintf("%s:%s", host, port))
}
