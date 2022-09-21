package handler

import (
	"github.com/gin-gonic/gin"
)

func BuzHandler() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.String(200, "🐝")
	}
	return gin.HandlerFunc(fn)
}
