package generic

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
)

func validateAndForwardEvents(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf config.Generic, events []interface{}) {

}

func PostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Info().Msg("post handler called")
	}
	return gin.HandlerFunc(fn)
}

func BatchPostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Info().Msg("batch post handler called")
	}
	return gin.HandlerFunc(fn)
}
