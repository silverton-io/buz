package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/constants"
)

func AdvancingCookie(conf config.Cookie) gin.HandlerFunc {
	return func(c *gin.Context) {
		identityCookieValue, _ := c.Cookie(conf.Name)
		if identityCookieValue != "" {
			c.SetSameSite(http.SameSiteLaxMode)
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
			identityCookieValue = uuid.New().String()
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie(
				conf.Name,
				identityCookieValue,
				60*60*24*conf.TtlDays,
				conf.Path,
				conf.Domain,
				conf.Secure,
				false,
			)
		}
		c.Set(constants.IDENTITY, identityCookieValue)
		c.Next()
	}
}
