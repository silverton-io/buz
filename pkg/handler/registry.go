// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/registry"
)

type CacheIndex struct {
	Count   int      `json:"count"`
	Schemas []string `json:"schemas"`
}

func CachePurgeHandler(r *registry.Registry) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Debug().Msg("🟡 schema cache purged")
		r.Cache.Clear()
	}
	return gin.HandlerFunc(fn)
}

func RegistrySchemaHandler(r *registry.Registry) gin.HandlerFunc {
	fn := func(c *gin.Context) {

	}
	return gin.HandlerFunc(fn)
}

// func RegistryIndexHandler(){} TODO! Fix this.

func CacheIndexHandler(r *registry.Registry) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var schemaKeys = make([]string, 0)
		iter := r.Cache.NewIterator()
		for {
			entry := iter.Next()
			if entry == nil {
				break
			}
			schemaKey := string(entry.Key)
			schemaKeys = append(schemaKeys, schemaKey)
		}
		index := CacheIndex{
			Count:   len(schemaKeys),
			Schemas: schemaKeys,
		}
		c.JSON(http.StatusOK, index)
	}
	return gin.HandlerFunc(fn)
}
