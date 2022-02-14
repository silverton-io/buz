package generic

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
)

func PostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Info().Msg("post handler called")
	}
	return gin.HandlerFunc(fn)
}

func BatchPostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Info().Msg("batch post handler called")
	}
	return gin.HandlerFunc(fn)
}
