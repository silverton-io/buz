package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/gosnowplow/pkg/response"
)

func Healthcheck(c *gin.Context) {
	c.JSON(200, response.Ok)
}
