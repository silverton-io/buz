package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/config"
)

func CORS(conf config.Cors) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", strings.Join(conf.AllowOrigin, ", "))
		c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(conf.AllowCredentials))
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", strings.Join(conf.AllowMethods, ", "))
		c.Header("Access-Control-Max-Age", strconv.Itoa(conf.MaxAge))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
