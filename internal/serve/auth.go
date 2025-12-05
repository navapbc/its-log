package serve

import (
	"net/http"

	"cms.hhs.gov/its-log/internal/config"
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
		var apiConfig config.ApiConfig
		err := viper.Unmarshal(&apiConfig)
		// TODO: Handle config failure
		if err != nil {
			panic(err)
		}

		api_key := c.GetHeader("x-api-key")
		for _, key := range apiConfig.Keys {
			if key.Kind == config.LOGGING_KEY_KIND &&
				len(api_key) > 32 &&
				api_key == key.Key {
				return
			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		}
	}
}
