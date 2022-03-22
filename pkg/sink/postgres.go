package sink

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type PostgresSink struct{}

func (s *PostgresSink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing postgres sink")
}

func (s *PostgresSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {}

func (s *PostgresSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
}

func (s *PostgresSink) Close() {
	log.Debug().Msg("closing postgres sink")
}
