package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/tele"
)

func StatsHandler(meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(200, meta)
	}
	return gin.HandlerFunc(fn)
}
