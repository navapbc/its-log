package endpoints

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/navapbc/its-log/internal/base"
	"github.com/navapbc/its-log/internal/constants"
	"github.com/navapbc/its-log/internal/types"
)

// its-log is intended to be used by a single application.
// For local testing, set `api_key` in the config file.
// For production deployments, set `ITSLOG_API_KEY` to a
// random value at least 48 bytes long. This is intended to be
// a shared, symmetric key between the client and its-log.
//
// python -c 'import secrets ; print(secrets.token_urlsafe(48))'
//
// or
//
// openssl rand -base64 48
//
// would likely do the trick.
//
// This middleware sets the AppId for use downstream

func AuthMiddleWare(permissions []types.PermissionType) gin.HandlerFunc {
	return func(c *gin.Context) {
		api_key := c.GetHeader("x-api-key")
		for _, key := range base.LiveKeys {
			keylen := len(api_key)
			doesContain := slices.Contains(permissions, key.Permission)
			if doesContain && keylen >= constants.MINIMUM_API_KEY_LENGTH {
				if api_key == key.Key {
					c.Set(constants.ITSLOG_KEYID, key.KeyId)
					c.Set(constants.ITSLOG_APPID, key.AppId)
					return
				}
			}
		}
		// Otherwise, fail.
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
