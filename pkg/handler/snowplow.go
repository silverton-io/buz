package handler

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/config"
	e "github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/request"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/sink"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func buildEnvelopesFromRequest(c *gin.Context, conf config.Snowplow, meta *tele.Meta) []e.Envelope {
	var events []e.Envelope
	if c.Request.Method == "POST" {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not read request body")
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := snowplow.BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), conf, meta)
			schema := spEvent.Schema()
			envelope := e.Envelope{
				EventProtocol: protocol.SNOWPLOW,
				EventSchema:   schema,
				Tstamp:        time.Now(),
				Ip:            *spEvent.User_ipaddress,
				Payload:       spEvent,
			}
			events = append(events, envelope)
		}
	} else {
		params := request.MapParams(c)
		spEvent := snowplow.BuildEventFromMappedParams(c, params, conf, meta)
		schema := spEvent.Schema()
		envelope := e.Envelope{
			EventProtocol: protocol.SNOWPLOW,
			EventSchema:   schema,
			Tstamp:        time.Now(),
			Ip:            *spEvent.User_ipaddress,
			Payload:       spEvent,
		}
		events = append(events, envelope)
	}
	return events
}

func SnowplowRedirectHandler(conf config.Snowplow, meta *tele.Meta, cache *cache.SchemaCache, sink sink.Sink) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := buildEnvelopesFromRequest(c, conf, meta)
		validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, cache)
		sink.BatchPublishValidAndInvalid(ctx, protocol.SNOWPLOW, validEvents, invalidEvents, meta)
		redirectUrl, _ := c.GetQuery("u")
		c.Redirect(302, redirectUrl)
	}
	return gin.HandlerFunc(fn)
}

func SnowplowDefaultHandler(conf config.Snowplow, meta *tele.Meta, cache *cache.SchemaCache, sink sink.Sink) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := buildEnvelopesFromRequest(c, conf, meta)
		validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, cache)
		sink.BatchPublishValidAndInvalid(ctx, protocol.SNOWPLOW, validEvents, invalidEvents, meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
