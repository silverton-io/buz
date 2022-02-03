package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/gosnowplow/pkg/config"
)

func AdvancingCookie(config config.AdvancingCookie) gin.HandlerFunc {
	return func(c *gin.Context) {
		identityCookieValue, _ := c.Cookie(config.CookieName)
		if identityCookieValue != "" {
			c.SetCookie(
				config.CookieName,
				identityCookieValue,
				60*60*24*config.CookieTtlDays,
				config.CookiePath,
				config.CookieDomain,
				config.UseSecureCookie,
				false,
			)
		} else {
			identityCookieValue := uuid.New()
			c.SetCookie(
				config.CookieName,
				identityCookieValue.String(),
				60*60*24*config.CookieTtlDays,
				config.CookiePath,
				config.CookieDomain,
				config.UseSecureCookie,
				false,
			)
		}
		c.Set("identity", identityCookieValue)
		c.Next()
	}
}
