package sink

import (
	"context"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/request"
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
	go request.PostEnvelopes(s.relayUrl, validEnvelopes)
}

func (s *RelaySink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {
	log.Error().Msg("BatchPublishInvalid is disabled for relay sink")
}

func (s *RelaySink) Close() {
	log.Debug().Msg("closing relay sink") // no-op
}
