package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// its-log is intended to be used by a single application.
// For local testing, set `api_key` in the config file.
// For production deployments, set `ITSLOG_API_KEY` to a
// random value at least 32 bytes long. This is intended to be
// a shared, symmetric key between the client and its-log.
//
// python -c 'import secrets ; print(secrets.token_urlsafe(32))'
//
// would likely do the trick.

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.GetHeader("x-api-key")

		if len(api_key) > 32 && api_key == viper.GetString("api_key") {
			return
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
