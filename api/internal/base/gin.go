package base

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetOrPanic(c *gin.Context, key string) string {
	v, exists := c.Get(key)
	if !exists {
		panic(fmt.Sprintf("could not get key from gin context: %s", key))
	}
	return v.(string)
}
