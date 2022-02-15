package generic

import (
	"context"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/event"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/tidwall/gjson"
)

func PostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		e, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs
		schemaName := gjson.GetBytes(e, conf.Payload.SchemaKey).String()
		eventPayload := gjson.GetBytes(e, conf.Payload.RootKey)
		isValid, validationError, _ := validateEvent(eventPayload, schemaName, cache)
		if isValid {
			forwarder.PublishValid(ctx, e)
		} else {
			iE := event.InvalidEvent{
				ValidationError: &validationError,
				Event:           eventPayload.Value().(map[string]interface{}),
			}
			forwarder.PublishInvalid(ctx, iE)
		}
	}
	return gin.HandlerFunc(fn)
}

func BatchPostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		log.Info().Msg("batch post handler called")
	}
	return gin.HandlerFunc(fn)
}
