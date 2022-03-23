package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type BlackholeSink struct {
	id   *uuid.UUID
	name string
}

func (s *BlackholeSink) Id() *uuid.UUID {
	return s.id
}

func (s *BlackholeSink) Initialize(conf config.Sink) {
	id := uuid.New()
	s.id, s.name = &id, conf.Name
}

func (s *BlackholeSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) {}

func (s *BlackholeSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {
}

func (s *BlackholeSink) Close() {
}
