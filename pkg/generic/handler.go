package generic

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	e "github.com/silverton-io/gosnowplow/pkg/event"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	f "github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/response"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"github.com/tidwall/gjson"
)

func bifurcateEvents(events []gjson.Result, cache *cache.SchemaCache, conf *config.Generic) (validEvents []interface{}, invalidEvents []interface{}) {
	var vEvents []interface{}
	var invEvents []interface{}
	for _, event := range events {
		payloadSchemaName := event.Get(conf.Payload.RootKey + "." + conf.Payload.SchemaKey).String()
		payloadData := event.Get(conf.Payload.RootKey + "." + conf.Payload.DataKey)
		isValid, validationError, _ := validateEvent(payloadData, payloadSchemaName, cache, conf)
		if isValid {
			vEvents = append(vEvents, event.Value())
		} else {
			invalidEvent := e.InvalidEvent{
				ValidationError: &validationError,
				Event:           event.Value(),
			}
			invEvents = append(invEvents, invalidEvent)
		}
	}
	return vEvents, invEvents
}

func PostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf *config.Generic, meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		var events []gjson.Result
		event := gjson.ParseBytes(reqBody)
		events = append(events, event)
		validEvents, invalidEvents := bifurcateEvents(events, cache, conf)
		f.BatchPublishValidAndInvalid(input.GENERIC_INPUT, forwarder, validEvents, invalidEvents, meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}

func BatchPostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, conf *config.Generic, meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		var rawEvents []interface{}
		var events []gjson.Result
		err := json.Unmarshal(reqBody, &rawEvents)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error when unmarshaling request body")
			// TODO! Decide whether or not to return something bad here
		}
		for _, rawEvent := range rawEvents {
			marshaledEvent, err := json.Marshal(rawEvent)
			if err != nil {
				log.Error().Stack().Err(err).Msg("error when marshaling event")
				// TODO! Decide whether or not to return something bad here
			} else {
				event := gjson.ParseBytes(marshaledEvent)
				events = append(events, event)
			}
		}
		validEvents, invalidEvents := bifurcateEvents(events, cache, conf)
		f.BatchPublishValidAndInvalid(input.GENERIC_INPUT, forwarder, validEvents, invalidEvents, meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
