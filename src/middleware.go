package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdvancingCookieConfig struct {
	cookieName      string
	useSecureCookie bool
	cookieTtlDays   int
	cookiePath      string
	cookieDomain    string
}

func AdvancingCookieMiddleware(config AdvancingCookieConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		identityCookieValue, _ := c.Cookie(config.cookieName)
		if identityCookieValue != "" {
			c.SetCookie(
				config.cookieName,
				identityCookieValue,
				60*60*24*config.cookieTtlDays,
				config.cookiePath,
				config.cookieDomain,
				config.useSecureCookie,
				false,
			)
		} else {
			identityCookieValue := uuid.New()
			c.SetCookie(
				config.cookieName,
				identityCookieValue.String(),
				60*60*24*config.cookieTtlDays,
				config.cookiePath,
				config.cookieDomain,
				config.useSecureCookie,
				false,
			)
		}
		c.Set("identity", identityCookieValue)
		c.Next()
	}
}

type CORSConfig struct {
	allowOrigin      []string
	allowCredentials bool
	allowMethods     []string
	maxAge           int
}

func CORSMiddleware(config CORSConfig) gin.HandlerFunc {
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
