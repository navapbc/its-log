package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/itslog"
)

func addTestingEndpoints(rG *gin.RouterGroup, ch_evt_out chan<- *itslog.Event) {
	// Test data generation
	auth_testV1 := rG.Group("/")
	permissions := []itslog.PermissionType{itslog.Test}
	auth_testV1.Use(AuthMiddleWare(permissions))
	auth_testV1.POST("log/:date", Event(ch_evt_out, itslog.Test))
}
