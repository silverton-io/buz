// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package util

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/constants"
)

func GetIdentityOrFallback(c *gin.Context, conf config.Middleware) string {
	identity := c.GetString(constants.IDENTITY)
	if identity == "" {
		return conf.Fallback
	}
	return identity
}
