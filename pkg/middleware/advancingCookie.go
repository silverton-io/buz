package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/gosnowplow/pkg/config"
)

func AdvancingCookie(conf config.Cookie) gin.HandlerFunc {
	return func(c *gin.Context) {
		identityCookieValue, _ := c.Cookie(conf.Name)
		if identityCookieValue != "" {
			c.SetCookie(
				conf.Name,
				identityCookieValue,
				60*60*24*conf.TtlDays,
				conf.Path,
				conf.Domain,
				conf.Secure,
				false,
			)
		} else {
			identityCookieValue := uuid.New()
			c.SetCookie(
				conf.Name,
				identityCookieValue.String(),
				60*60*24*conf.TtlDays,
				conf.Path,
				conf.Domain,
				conf.Secure,
				false,
			)
		}
		c.Set("identity", identityCookieValue)
		c.Next()
	}
}
