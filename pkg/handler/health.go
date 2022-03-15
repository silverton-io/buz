package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/response"
)

func HealthcheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, response.Ok)
}
