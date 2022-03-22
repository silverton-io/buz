package sink

import (
	"context"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/request"
)

type HttpSink struct {
	validUrl   url.URL
	invalidUrl url.URL
}

func (s *HttpSink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing http sink")
	vUrl, vErr := url.Parse(conf.ValidUrl)
	invUrl, invErr := url.Parse(conf.InvalidUrl)
	if vErr != nil {
		log.Fatal().Stack().Err(vErr).Msg("validUrl is not a valid url")
	}
	if invErr != nil {
		log.Fatal().Stack().Err(invErr).Msg("invalidUrl is not a valid url")
	}
	s.validUrl = *vUrl
	s.invalidUrl = *invUrl
}

func (s *HttpSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) {
	request.PostEnvelopes(s.validUrl, validEnvelopes)
}

func (s *HttpSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {
	request.PostEnvelopes(s.invalidUrl, invalidEnvelopes)
}

func (s *HttpSink) Close() {
	log.Debug().Msg("closing http sink") // no-op
}
