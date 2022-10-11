// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
)

func Identity(conf config.Identity) gin.HandlerFunc {
	return func(c *gin.Context) {
		identityCookieValue, _ := c.Cookie(conf.Cookie.Name)
		switch conf.Cookie.SameSite {
		case "None":
			c.SetSameSite(http.SameSiteNoneMode)
		case "Lax":
			c.SetSameSite(http.SameSiteLaxMode)
		case "Strict":
			c.SetSameSite(http.SameSiteStrictMode)
		}
		if identityCookieValue != "" {
			c.SetCookie(
				conf.Cookie.Name,
				identityCookieValue,
				60*60*24*conf.Cookie.TtlDays,
				conf.Cookie.Path,
				conf.Cookie.Domain,
				conf.Cookie.Secure,
				false,
			)
		} else {
			identityCookieValue = uuid.New().String()
			c.SetCookie(
				conf.Cookie.Name,
				identityCookieValue,
				60*60*24*conf.Cookie.TtlDays,
				conf.Cookie.Path,
				conf.Cookie.Domain,
				conf.Cookie.Secure,
				false,
			)
		}
		c.Set(constants.IDENTITY, identityCookieValue)
		c.Next()
	}
}
