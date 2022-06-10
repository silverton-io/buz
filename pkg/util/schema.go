package util

import (
	"github.com/gin-gonic/gin"
)

func GetSchemaNameFromRequest(c *gin.Context, fallback string) string {
	schemaName := c.Param("hps")
	if schemaName == "" || schemaName == "/" { // Handle no param, or trailing slash
		return fallback
	}
	schemaName = schemaName[1:] + JSON_EXTENSION
	return schemaName
}
