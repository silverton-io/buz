package sink

import (
	"context"

	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type BlackholeSink struct{}

func (s *BlackholeSink) Initialize(conf config.Sink) {}

func (s *BlackholeSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) {}

func (s *BlackholeSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {
}

func (s *BlackholeSink) Close() {
}
