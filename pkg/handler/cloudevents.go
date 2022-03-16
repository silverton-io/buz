package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func CloudeventsPostHandler(handlerParams EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var envelopes []event.Envelope
		if c.ContentType() == "application/cloudevents+json" { // Only accept request if content type is set appropriately
			ctx := context.Background()
			reqBody, _ := ioutil.ReadAll(c.Request.Body)
			result := gjson.ParseBytes(reqBody)
			cloudevent := cloudevents.BuildEvent(c, result.Value().(map[string]interface{}))
			envelope := event.Envelope{
				EventProtocol: protocol.CLOUDEVENTS,
				EventSchema:   &cloudevent.DataSchema,
				Tstamp:        time.Now(),
				Payload:       cloudevent,
			}
			envelopes = append(envelopes, envelope)
			validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, handlerParams.Cache)
			handlerParams.Sink.BatchPublishValidAndInvalid(ctx, protocol.CLOUDEVENTS, validEvents, invalidEvents, handlerParams.Meta)
			c.JSON(200, response.Ok)
		} else {
			c.JSON(400, response.InvalidContentType)
		}

	}
	return gin.HandlerFunc(fn)
}

func CloudeventsBatchPostHandler(handlerParams EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		var envelopes []event.Envelope
		if c.ContentType() == "application/cloudevents-batch+json" {
			reqBody, _ := ioutil.ReadAll(c.Request.Body)
			var rawEvents []interface{}
			err := json.Unmarshal(reqBody, &rawEvents)
			if err != nil {
				log.Error().Stack().Err(err).Msg("error when unmarshaling request body")
				c.JSON(400, response.BadRequest)
				return
			}
			for _, rawEvent := range rawEvents {
				marshaledEvent, err := json.Marshal(rawEvent)
				if err != nil {
					log.Error().Stack().Err(err).Msg("error when marshaling event")
					c.JSON(400, response.BadRequest)
					return
				} else {
					e := gjson.ParseBytes(marshaledEvent)
					cloudevent := cloudevents.BuildEvent(c, e.Value().(map[string]interface{}))
					envelope := event.Envelope{
						EventProtocol: protocol.CLOUDEVENTS,
						EventSchema:   &cloudevent.DataSchema,
						Tstamp:        time.Now(),
						Payload:       cloudevent,
					}
					envelopes = append(envelopes, envelope)
				}
			}
			validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, handlerParams.Cache)
			handlerParams.Sink.BatchPublishValidAndInvalid(ctx, protocol.CLOUDEVENTS, validEvents, invalidEvents, handlerParams.Meta)
			c.JSON(http.StatusOK, response.Ok)
		} else {
			c.JSON(http.StatusBadRequest, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)
}
