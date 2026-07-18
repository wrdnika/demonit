package http

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	HeaderDeviceAPIKey = "X-Device-API-Key"
	HeaderAdminAPIKey  = "X-Admin-API-Key"
)

// requireAPIKey validates a shared secret header using constant-time comparison.
func requireAPIKey(headerName, expected string) gin.HandlerFunc {
	return func(c *gin.Context) {
		provided := c.GetHeader(headerName)
		if expected == "" || subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) != 1 {
			JSONError(c, http.StatusUnauthorized, "UNAUTHORIZED", "invalid or missing API key", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
