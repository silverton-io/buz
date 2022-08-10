// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	cache "github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/response"
)

type CacheIndex struct {
	Count   int      `json:"count"`
	Schemas []string `json:"schemas"`
}

func CachePurgeHandler(s *cache.SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Debug().Msg("schema cache purged")
		s.Cache.Clear()
	}
	return gin.HandlerFunc(fn)
}

func CacheIndexHandler(s *cache.SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var schemaKeys = make([]string, 0)
		iter := s.Cache.NewIterator()
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

func CacheGetHandler(s *cache.SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		schemaName := c.Param(cache.SCHEMA_ROUTE_PARAM)[1:]
		cacheKey := []byte(schemaName)
		cachedSchema, _ := s.Cache.Get(cacheKey)
		if cachedSchema != nil {
			var schema interface{}
			err := json.Unmarshal(cachedSchema, &schema)
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not unmarshal cached schema")
				c.JSON(http.StatusBadRequest, response.BadRequest)
				return
			}
			c.JSON(http.StatusOK, schema)
			return
		} else {
			c.JSON(http.StatusNotFound, response.SchemaNotCached)
			return
		}
	}
	return gin.HandlerFunc(fn)
}
