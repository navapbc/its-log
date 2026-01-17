package serve

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/jadudm/its-log/internal/itslog"
)

// its-log is intended to be used by a single application.
// For local testing, set `api_key` in the config file.
// For production deployments, set `ITSLOG_API_KEY` to a
// random value at least 32 bytes long. This is intended to be
// a shared, symmetric key between the client and its-log.
//
// python -c 'import secrets ; print(secrets.token_urlsafe(32))'
//
// or
//
// openssl rand -hex 32
//
// would likely do the trick.
//
// This middleware sets the AppId for use downstream

func AuthMiddleWare(permissions []itslog.PermissionType) gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.GetHeader("x-api-key")
		for _, key := range itslog.LiveKeys {
			keylen := len(api_key)
			doesContain := slices.Contains(permissions, key.Permission)
			if doesContain && keylen >= 32 {
				if api_key == key.Key {
					c.Set("KeyId", key.KeyId)
					c.Set("AppId", key.AppId)
					return
				}
			}
		}
		// Otherwise, fail.
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
