package envelope

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/tidwall/gjson"
)

func BuildGenericEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		n := buildCommonEnvelope(c, conf.Middleware, m)
		genEvent, err := generic.BuildEvent(e, conf.Generic)
		if err != nil {
			log.Error().Err(err).Msg("could not build generic event")
		}
		// Event meta
		n.EventMeta.Protocol = protocol.GENERIC
		n.EventMeta.Schema = genEvent.Payload.Schema
		// Context
		n.Contexts = genEvent.Contexts
		// Payload
		n.Payload = genEvent.Payload.Data
		envelopes = append(envelopes, n)
	}
	return envelopes
}
