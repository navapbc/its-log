package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
)

func addTestingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *base.Event) {
	// Test data generation
	auth_testV1 := rG.Group("/")
	permissions := []base.PermissionType{base.Test}
	auth_testV1.Use(AuthMiddleWare(permissions))
	auth_testV1.POST("log/:date", Event(ch_evt_out, base.Test))
}
