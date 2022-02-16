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
	"github.com/silverton-io/gosnowplow/pkg/util"
	"github.com/tidwall/gjson"
)

func bifurcateEvents(events []interface{}, cache *cache.SchemaCache, conf *config.Generic) (validEvents []interface{}, invalidEvents []interface{}) {
	var vEvents []interface{}
	var invEvents []interface{}
	for _, event := range events {
		marshaledEvent, _ := json.Marshal(event)
		gResult := gjson.ParseBytes(marshaledEvent)
		util.PrettyPrint(gResult.Value())
		payloadSchemaName := gResult.Get(conf.Payload.RootKey + "." + conf.Payload.SchemaKey).String()
		payloadData := gResult.Get(conf.Payload.RootKey + "." + conf.Payload.DataKey)
		util.PrettyPrint(payloadSchemaName)
		util.PrettyPrint(payloadData)
		isValid, validationError, _ := validateEvent(payloadData, payloadSchemaName, cache, conf)
		if isValid {
			vEvents = append(vEvents, gResult.Value())
		} else {
			invalidEvent := e.InvalidEvent{
				ValidationError: &validationError,
				Event:           gResult.Value(),
			}
			invEvents = append(invEvents, invalidEvent)
		}
	}
	return vEvents, invEvents
}

func PostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf *config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		reqBody, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs
		var events []interface{}
		e := gjson.ParseBytes(reqBody)
		events = append(events, e)
		validEvents, invalidEvents := bifurcateEvents(events, cache, conf)
		forwarder.BatchPublishValid(ctx, validEvents)
		forwarder.BatchPublishInvalid(ctx, invalidEvents)
	}
	return gin.HandlerFunc(fn)
}

func BatchPostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf *config.Generic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		var events []interface{}
		json.Unmarshal(reqBody, &events)
		validEvents, invalidEvents := bifurcateEvents(events, cache, conf)
		forwarder.BatchPublishValid(ctx, validEvents)
		forwarder.BatchPublishInvalid(ctx, invalidEvents)
	}
	return gin.HandlerFunc(fn)
}
