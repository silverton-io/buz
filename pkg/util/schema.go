package util

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/constants"
)

func GetSchemaNameFromRequest(c *gin.Context, fallback string) string {
	schemaName := c.Param(constants.HONEYPOT_SCHEMA_PARAM)
	if schemaName == "" || schemaName == "/" { // Handle no param, or trailing slash
		return fallback
	}
	schemaName = schemaName[1:] + JSON_EXTENSION
	return schemaName
}
