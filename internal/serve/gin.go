package serve

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/itslog"
	"github.com/spf13/viper"
)

func PourGin(s itslog.ItsLog) *gin.Engine {
	// We may want production mode.
	// This is configured via the envrionment
	if viper.GetString("gin_mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	apiV1 := router.Group("/v1")
	authV1 := apiV1.Group("/")
	authV1.Use(AuthMiddleWare())

	authV1.POST("log", LogV1(s))
	return router
}

func LogV1(s itslog.ItsLog) func(c *gin.Context) {
	return func(c *gin.Context) {

		// The middleware already did the binding check.
		// If it failed, we never got here. So, in this function,
		// we can ignore the error handling. If that ever becomes
		// a performance issue, the middleware can be removed, as it
		// represents a double-binding/allocation
		event, err := BindEvent(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "err",
				"err":    err.Error(),
			})
		}
		id, err := s.Event(event.Source, event.Event, event.Value, event.Type)

		if err != nil {
			c.JSON(http.StatusTeapot, gin.H{
				"status": "err",
				"err":    err.Error(),
			})
		}

		// Everything worked.
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"id":     id,
		})

	}

}
