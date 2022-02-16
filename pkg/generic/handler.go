package generic

import (
	"context"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/event"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/util"
	"github.com/tidwall/gjson"
)

func PostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		e, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs

		jsonEvent := gjson.ParseBytes(e)

		schemaName := jsonEvent.Get(conf.Payload.RootKey + "." + conf.Payload.SchemaKey).String()
		eventPayload := jsonEvent.Get(conf.Payload.RootKey + "." + conf.Payload.DataKey)

		isValid, validationError, _ := validateEvent(eventPayload, schemaName, cache)
		if isValid {
			forwarder.PublishValid(ctx, jsonEvent.Value())
			util.PrettyPrint(jsonEvent.Value())
		} else {
			iE := event.InvalidEvent{
				ValidationError: &validationError,
				Event:           jsonEvent.Value(),
			}
			util.PrettyPrint(iE)
			forwarder.PublishInvalid(ctx, iE)
		}
	}
	return gin.HandlerFunc(fn)
}

func BatchPostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// ctx := context.Background()
		events, _ := ioutil.ReadAll(c.Request.Body)
		jsonEvents := gjson.GetManyBytes(events)
		for _, r := range jsonEvents {
			util.PrettyPrint(r.Value())
		}
		util.PrettyPrint(jsonEvents)
	}
	return gin.HandlerFunc(fn)
}
