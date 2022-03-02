package cache

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func CachePurgeHandler(s *SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Debug().Msg("schema cache purged")
		s.cache.Clear()
	}
	return gin.HandlerFunc(fn)
}
