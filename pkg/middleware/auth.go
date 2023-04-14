// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/response"
)

const (
	BASIC  = "Basic"
	BEARER = "Bearer"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

func tokenIsValid(token string, conf config.Auth) bool {
	for _, validToken := range conf.Tokens {
		if token == validToken {
			return true
		}
	}
	return false
}

// The simplest-possible way to lock down routes
func Auth(conf config.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}
		if err := c.ShouldBindHeader(&h); err != nil {
			// Can't bind Authorization header
			c.JSON(http.StatusUnauthorized, response.MissingAuthHeader)
			c.Abort()
			return
		}
		if h.Token == "" {
			// No Authorization header present
			c.JSON(http.StatusUnauthorized, response.MissingAuthHeader)
			c.Abort()
			return
		}
		tokenParts := strings.Split(h.Token, " ")
		if len(tokenParts) < 2 {
			// Header present but missing scheme or token
			c.JSON(http.StatusUnauthorized, response.MissingAuthSchemeOrToken)
			c.Abort()
			return
		}
		scheme := tokenParts[0]
		token := tokenParts[1]
		if scheme != BEARER && scheme != BASIC {
			// Auth scheme isn't supported
			c.JSON(http.StatusUnauthorized, response.InvalidAuthScheme)
			c.Abort()
			return
		}
		isValid := tokenIsValid(token, conf)
		if !isValid {
			// Invalid token
			c.JSON(http.StatusUnauthorized, response.InvalidAuthToken)
			c.Abort()
			return
		}
		c.Next()
	}
}
