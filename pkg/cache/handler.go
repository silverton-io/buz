package cache

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/response"
)

type CacheIndex struct {
	Count   int      `json:"count"`
	Schemas []string `json:"schemas"`
}

func CachePurgeHandler(s *SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Debug().Msg("schema cache purged")
		s.cache.Clear()
	}
	return gin.HandlerFunc(fn)
}

func CacheIndexHandler(s *SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var schemaKeys = make([]string, 0)
		iter := s.cache.NewIterator()
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

func CacheGetHandler(s *SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		schemaName := c.Param("schema")[1:]
		cacheKey := []byte(schemaName)
		cachedSchema, _ := s.cache.Get(cacheKey)
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
