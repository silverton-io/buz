// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/constants"
)

func GetSchemaNameFromRequest(c *gin.Context, fallback string) string {
	schemaName := c.Param(constants.BUZ_SCHEMA_PARAM)
	if schemaName == "" || schemaName == "/" { // Handle no param, or trailing slash
		return fallback
	}
	schemaName = schemaName[1:] + JSON_EXTENSION
	return schemaName
}
