package sink

import (
	"context"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/request"
	"github.com/silverton-io/honeypot/pkg/tele"
)

type RelaySink struct {
	relayUrl url.URL
}

func (s *RelaySink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing http sink")
	u, err := url.Parse(conf.RelayUrl)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("relayUrl is not a valid url")
	}
	s.relayUrl = *u
}

func (s *RelaySink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) {
	log.Error().Msg("BatchPublishValid is disabled for relay sink")
}

func (s *RelaySink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {
	log.Error().Msg("BatchPublishInvalid is disabled for relay sink")
}

func (s *RelaySink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope, meta *tele.Meta) {
	var envelopes []envelope.Envelope
	envelopes = append(validEnvelopes, invalidEnvelopes...)
	go request.PostEnvelopes(s.relayUrl, envelopes)
	// FIXME! Increment stats. Not including this yet because want to go event protocol/name/etc route.
}

func (s *RelaySink) Close() {
	log.Debug().Msg("closing http sink") // no-op
}
