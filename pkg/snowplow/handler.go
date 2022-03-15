package snowplow

import (
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
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func buildEventsFromRequest(c *gin.Context, conf config.Snowplow, meta *tele.Meta) []e.Envelope {
	var events []e.Envelope
	if c.Request.Method == "POST" {
		body, err := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs here
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not read request body")
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), conf, meta)
			envelope := e.Envelope{
				EventProtocol: protocol.SNOWPLOW,
				// EventSchema:   &spEvent.Self_describing_event.Schema,
				Tstamp: time.Now(),
				Ip:     *spEvent.User_ipaddress,
				Event:  spEvent.toMap(),
			}
			events = append(events, envelope)
		}
	} else {
		params := request.MapParams(c)
		spEvent := BuildEventFromMappedParams(c, params, conf, meta)
		envelope := e.Envelope{
			EventProtocol: protocol.SNOWPLOW,
			// EventSchema:   &spEvent.Self_describing_event.Schema,
			Tstamp: time.Now(),
			Ip:     *spEvent.User_ipaddress,
			Event:  spEvent.toMap(),
		}
		events = append(events, envelope)
	}
	return events
}

func RedirectHandler(conf config.Snowplow, meta *tele.Meta, cache *cache.SchemaCache, sink sink.Sink) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		// ctx := context.Background()
		// events := buildEventsFromRequest(c, conf, meta)
		// validEvents, invalidEvents := bifurcateEvents(events, cache)
		// sink.BatchPublishValidAndInvalid(ctx, protocol.SNOWPLOW, validEvents, invalidEvents, meta)
		redirectUrl, _ := c.GetQuery("u")
		c.Redirect(302, redirectUrl)
	}
	return gin.HandlerFunc(fn)
}

func DefaultHandler(conf config.Snowplow, meta *tele.Meta, cache *cache.SchemaCache, sink sink.Sink) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		v := validator.Validator{}
		// ctx := context.Background()
		eventEnvelopes := buildEventsFromRequest(c, conf, meta)
		for _, envelope := range eventEnvelopes {
			v.ValidateEnvelope(&envelope)
			util.Pprint(envelope)
		}
		// validEvents, invalidEvents := bifurcateEvents(events, cache)
		// sink.BatchPublishValidAndInvalid(ctx, protocol.SNOWPLOW, validEvents, invalidEvents, meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
