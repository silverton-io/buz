package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AdvancingCookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		identityCookieValue, _ := c.Cookie(Config.cookieName)
		var useSecureCookie bool
		if Config.environment == "dev" {
			useSecureCookie = false
		} else {
			useSecureCookie = true
		}
		if identityCookieValue != "" {
			c.SetCookie(
				Config.cookieName,
				identityCookieValue,
				60*60*24*Config.cookieTtlDays,
				Config.cookiePath,
				Config.cookieDomain,
				useSecureCookie,
				false,
			)
		} else {
			identityCookieValue := uuid.New()
			c.SetCookie(
				Config.cookieName,
				identityCookieValue.String(),
				60*60*24*Config.cookieTtlDays,
				Config.cookiePath,
				Config.cookieDomain,
				useSecureCookie,
				false,
			)
		}
		c.Set("identity", identityCookieValue)
		c.Next()
	}
}
