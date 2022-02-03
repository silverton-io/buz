package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/gosnowplow/pkg/config"
)

func CORS(config config.Cors) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // FIXME! Make this configurable.
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET") // FIXME! Make this configurable.
		c.Header("Access-Control-Max-Age", "86400")                    // FIXME! Make this configurable.

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
