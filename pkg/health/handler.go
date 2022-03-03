package health

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/response"
)

func HealthcheckHandler(c *gin.Context) {
	c.JSON(200, response.Ok)
}
