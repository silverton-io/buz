// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/registry"
	"github.com/silverton-io/buz/pkg/response"
	"github.com/tidwall/gjson"
)

func RegistryCachePurgeHandler(r *registry.Registry) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Debug().Msg("🟡 schema cache purged")
		r.Cache.Clear()
	}
	return gin.HandlerFunc(fn)
}

func RegistryGetSchemaHandler(r *registry.Registry) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		schemaName := c.Param(registry.SCHEMA_ROUTE_PARAM)[1:]
		exists, schemaContents := r.Get(schemaName)
		if !exists {
			c.JSON(404, response.SchemaNotAvailable)
		} else {
			schema := gjson.ParseBytes(schemaContents)
			c.JSON(200, schema.Value())
		}

	}
	return gin.HandlerFunc(fn)
}
