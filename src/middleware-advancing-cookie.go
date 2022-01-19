package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AdvancingCookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the tracking cookie if it doesn't already exist
		// If it does exist, set TTL "N time from request"
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
			cookieValue := uuid.New()
			c.SetCookie(
				Config.cookieName,
				cookieValue.String(),
				60*60*24*Config.cookieTtlDays,
				Config.cookiePath,
				Config.cookieDomain,
				useSecureCookie,
				false,
			)
		}
		c.Next()
	}
}
