// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
)

func GetIdentityOrFallback(c *gin.Context, conf config.Identity) string {
	identity := c.GetString(constants.IDENTITY)
	if identity == "" {
		return conf.Fallback
	}
	return identity
}
