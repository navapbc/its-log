package serve

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/config"
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

func AuthMiddleWare(auth_kind string, apiKeys config.ApiKeys) gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.GetHeader("x-api-key")
		for _, key := range apiKeys {
			if key.Kind == auth_kind && len(api_key) > 32 {
				log.Printf("%s %s\n", auth_kind, key.Key)
				// If the key is the right kind
				if api_key == key.Key {
					return
				}
			}
		}
		// Otherwise, fail.
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
