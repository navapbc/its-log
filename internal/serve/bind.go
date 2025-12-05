package serve

import (
	"cms.hhs.gov/its-log/internal/itslog"
	"github.com/gin-gonic/gin"
)

func BindEvent(c *gin.Context) (*itslog.Event, error) {
	var event itslog.Event
	err := c.BindJSON(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
