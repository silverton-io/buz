package handler

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/sink"
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/tidwall/gjson"
)

// func bifurcateEvents(e []event.CloudEvent, cache *cache.SchemaCache) (validEvents []interface{}, invalidEvents []interface{}) {
// 	var vEvents []interface{}
// 	var invEvents []interface{}
// 	for _, event := range events {
// 		isValid, validationError := validateEvent(event, cache)
// 		if isValid {
// 			vEvents = append(vEvents, event)
// 		} else {
// 			invalidEvent := e.InvalidEvent{
// 				ValidationError: &validationError,
// 				Event:           event,
// 			}
// 			invEvents = append(invEvents, invalidEvent)
// 		}
// 	}
// 	return vEvents, invEvents
// }

func CloudeventsPostHandler(conf *config.Cloudevents, meta *tele.Meta, cache *cache.SchemaCache, sink sink.Sink) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var events []event.Envelope
		if c.ContentType() == "application/cloudevents+json" { // Only accept request if content type is set appropriately
			// ctx := context.Background()
			reqBody, _ := ioutil.ReadAll(c.Request.Body)
			result := gjson.ParseBytes(reqBody)
			cloudevent := cloudevents.BuildEvent(c, result.Value().(map[string]interface{}))
			envelope := event.Envelope{
				EventProtocol: protocol.CLOUDEVENTS,
				EventSchema:   &cloudevent.DataSchema,
				Tstamp:        time.Now(),
				Payload:       cloudevent,
			}
			events = append(events, envelope)
			util.Pprint(events)
			// validEvents, invalidEvents := bifurcateEvents(cloudevents, cache)
			// sink.BatchPublishValidAndInvalid(ctx, protocol.CLOUDEVENTS, validEvents, invalidEvents, meta)
			c.JSON(200, response.Ok)
		} else {
			c.JSON(400, response.InvalidContentType)
		}

	}
	return gin.HandlerFunc(fn)
}

func CloudeventsBatchPostHandler(conf *config.Cloudevents, meta *tele.Meta, cache *cache.SchemaCache, sink sink.Sink) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// ctx := context.Background()
		// var events []cloudevents.CloudEvent
		if c.ContentType() == "application/cloudevents-batch+json" {
			reqBody, _ := ioutil.ReadAll(c.Request.Body)
			var rawEvents []interface{}
			err := json.Unmarshal(reqBody, &rawEvents)
			if err != nil {
				log.Error().Stack().Err(err).Msg("error when unmarshaling request body")
				c.JSON(400, response.BadRequest)
				return
			}
			// for _, rawEvent := range rawEvents {
			// 	if err != nil {
			// 		log.Error().Stack().Err(err).Msg("error when marshaling event")
			// 		c.JSON(400, response.BadRequest)
			// 		return
			// 	} else {
			// 		event := gjson.ParseBytes(marshaledEvent)
			// 		cloudevent := buildCloudevent(event)
			// 		cloudevents = append(cloudevents, cloudevent)
			// 	}
			// }
			// validEvents, invalidEvents := bifurcateEvents(cloudevents, cache)
			// sink.BatchPublishValidAndInvalid(ctx, protocol.CLOUDEVENTS, validEvents, invalidEvents, meta)
			c.JSON(200, response.Ok)
		} else {
			c.JSON(400, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)
}
