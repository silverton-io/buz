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
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func buildGenericEnvelopesFromRequest(c *gin.Context, conf config.Config) []event.Envelope {
	var envelopes []event.Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
	}
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		genEvent := generic.BuildEvent(e, conf.Generic)
		envelope := event.Envelope{
			EventProtocol: protocol.GENERIC,
			EventSchema:   &genEvent.Payload.Schema,
			Tstamp:        time.Now(),
			Ip:            c.ClientIP(),
			Payload:       genEvent,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func GenericHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := buildGenericEnvelopesFromRequest(c, *p.Config)
		validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, p.Cache)
		p.Sink.BatchPublishValidAndInvalid(ctx, protocol.GENERIC, validEvents, invalidEvents, p.Meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
