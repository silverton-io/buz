package envelope

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/tidwall/gjson"
)

func BuildCloudeventEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	for _, ce := range gjson.ParseBytes(reqBody).Array() {
		n := buildCommonEnvelope(c, m)
		cEvent, err := cloudevents.BuildEvent(ce)
		if err != nil {
			log.Error().Err(err).Msg("could not build Cloudevent")
		}
		// Event Meta
		n.EventMeta.Protocol = protocol.CLOUDEVENTS
		// Source
		n.Pipeline.Source.GeneratedTstamp = cEvent.Time
		n.Pipeline.Source.SentTstamp = cEvent.Time
		// Payload
		n.Payload = cEvent
		envelopes = append(envelopes, n)
	}
	return envelopes
}
