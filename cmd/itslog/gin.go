package main

import (
	"net/http"

	"cms.hhs.gov/its-log/internal/itslog"
	"github.com/gin-gonic/gin"
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
	authV1 := apiV1.Group("/", AuthMiddleWare())

	authV1.POST("log", LogV1(s))
	return router
}

func LogV1(s itslog.ItsLog) func(c *gin.Context) {
	return func(c *gin.Context) {
		var event itslog.Event
		if err := c.BindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "err",
				"error":  err.Error()})
			return
		}

		id, err := s.Event(event.Event, event.Value, itslog.LogItType(event.ValueType))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "err",
				"id":     err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"id":     id,
		})

	}

}
