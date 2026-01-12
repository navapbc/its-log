package serve

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetKeyId(c *gin.Context) (string, error) {
	key_id_any, exists := c.Get("key_id")
	if !exists {
		return "", errors.New("could not get key id from Gin context")
	}

	key_id := key_id_any.(string)
	if len(key_id) < 1 {
		return "", errors.New("key_id length too short")
	}

	return key_id, nil
}
