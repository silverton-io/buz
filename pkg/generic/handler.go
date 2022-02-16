package generic

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	e "github.com/silverton-io/gosnowplow/pkg/event"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/tidwall/gjson"
)

// FIXME! Abstract shared logic
func PostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf *config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		reqBody, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs

		event := gjson.ParseBytes(reqBody)
		payloadSchemaName := event.Get(conf.Payload.RootKey + "." + conf.Payload.SchemaKey).String()
		payloadData := event.Get(conf.Payload.RootKey + "." + conf.Payload.DataKey)

		isValid, validationError, _ := validateEvent(payloadData, payloadSchemaName, cache)
		if isValid {
			forwarder.PublishValid(ctx, event.Value())
		} else {
			invalidEvent := e.InvalidEvent{
				ValidationError: &validationError,
				Event:           event.Value(),
			}
			forwarder.PublishInvalid(ctx, invalidEvent)
		}
	}
	return gin.HandlerFunc(fn)
}

// FIXME! Abstract shared logic
func BatchPostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf *config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		var validEvents []interface{}
		var invalidEvents []interface{}
		var events []interface{}
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(reqBody, &events)

		for _, evnt := range events {
			marshaledEvent, _ := json.Marshal(evnt) // It feels incredibly dirty to json.Unmarshal a request body, only to json.Marshal and then gjson.ParseBytes.
			event := gjson.ParseBytes(marshaledEvent)
			payloadSchemaName := event.Get(conf.Payload.RootKey + "." + conf.Payload.SchemaKey).String()
			payloadData := event.Get(conf.Payload.RootKey + "." + conf.Payload.DataKey)
			isValid, validationError, _ := validateEvent(payloadData, payloadSchemaName, cache)
			if isValid {
				validEvents = append(validEvents, event)
			} else {
				invalidEvent := e.InvalidEvent{
					ValidationError: &validationError,
					Event:           event.Value(),
				}
				invalidEvents = append(invalidEvents, invalidEvent)
			}
		}
		forwarder.BatchPublishValid(ctx, validEvents)
		forwarder.BatchPublishInvalid(ctx, invalidEvents)
	}
	return gin.HandlerFunc(fn)
}
