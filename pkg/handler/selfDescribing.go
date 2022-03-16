package handler

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func buildGenericEnvelopesFromRequest(c *gin.Context, conf config.Config) []event.Envelope {
	var envelopes []event.Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
	}
	// Handles case of single-event
	e := gjson.ParseBytes(reqBody)
	util.Pprint(e.Value())
	genEvent := generic.BuildEvent(e, conf.Generic)
	envelope := event.Envelope{
		EventProtocol: protocol.GENERIC,
		EventSchema:   &genEvent.Payload.Schema,
		Tstamp:        time.Now(),
		Ip:            c.ClientIP(),
		Payload:       genEvent,
	}
	// TODO - handle case of batched events
	envelopes = append(envelopes, envelope)
	return envelopes
}

func GenericPostHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := buildGenericEnvelopesFromRequest(c, *p.Config)
		validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, p.Cache)
		p.Sink.BatchPublishValidAndInvalid(ctx, protocol.GENERIC, validEvents, invalidEvents, p.Meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}

func GenericBatchPostHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// ctx := context.Background()

		// var rawEvents []interface{}

		// err := json.Unmarshal(reqBody, &rawEvents)
		// if err != nil {
		// 	log.Error().Stack().Err(err).Msg("error when unmarshaling request body")
		// 	// TODO! Decide whether or not to return something bad here
		// }
		buildGenericEnvelopesFromRequest(c, *p.Config)

		// for _, rawEvent := range rawEvents {
		// 	marshaledEvent, err := json.Marshal(rawEvent)
		// 	if err != nil {
		// 		log.Error().Stack().Err(err).Msg("error when marshaling event")
		// 		// TODO! Decide whether or not to return something bad here
		// 	} else {
		// 		event := gjson.ParseBytes(marshaledEvent)
		// 		events = append(events, event)
		// 	}
		// }
		// validEvents, invalidEvents := validator.BifurcateAndAnnotate(events, p.Cache)
		// p.Sink.BatchPublishValidAndInvalid(ctx, protocol.GENERIC, validEvents, invalidEvents, p.Meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
