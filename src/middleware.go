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

func CORSMiddleware() gin.HandlerFunc {
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
