package serve

import (
	"github.com/gin-gonic/gin"
	itslog "github.com/jadudm/its-log/internal/itslog"
)

func BindEvent(c *gin.Context) (*itslog.Event, error) {
	var event itslog.Event
	err := c.BindJSON(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
