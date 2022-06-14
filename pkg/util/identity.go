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
