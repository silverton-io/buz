package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	e "github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/request"
	"github.com/silverton-io/honeypot/pkg/response"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/validator"
	"github.com/tidwall/gjson"
)

func buildEnvelopesFromRequest(c *gin.Context, conf config.Config) []e.Envelope {
	var events []e.Envelope
	if c.Request.Method == "POST" {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not read request body")
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := snowplow.BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), conf)
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
		spEvent := snowplow.BuildEventFromMappedParams(c, params, conf)
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

func SnowplowHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		envelopes := buildEnvelopesFromRequest(c, *p.Config)
		validEvents, invalidEvents := validator.BifurcateAndAnnotate(envelopes, p.Cache)
		p.Sink.BatchPublishValidAndInvalid(ctx, protocol.SNOWPLOW, validEvents, invalidEvents, p.Meta)
		if c.Request.Method == http.MethodGet {
			redirectUrl, _ := c.GetQuery("u")
			if redirectUrl != "" && p.Config.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("redirecting to " + redirectUrl)
				c.Redirect(http.StatusFound, redirectUrl)
			}
		}
		c.JSON(http.StatusOK, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
