package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/gosnowplow/pkg/config"
)

func AdvancingCookie(config config.Cookie) gin.HandlerFunc {
	return func(c *gin.Context) {
		identityCookieValue, _ := c.Cookie(config.Name)
		if identityCookieValue != "" {
			c.SetCookie(
				config.Name,
				identityCookieValue,
				60*60*24*config.TtlDays,
				config.Path,
				config.Domain,
				config.Secure,
				false,
			)
		} else {
			identityCookieValue := uuid.New()
			c.SetCookie(
				config.Name,
				identityCookieValue.String(),
				60*60*24*config.TtlDays,
				config.Path,
				config.Domain,
				config.Secure,
				false,
			)
		}
		c.Set("identity", identityCookieValue)
		c.Next()
	}
}
