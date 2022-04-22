package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type BlackholeSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
}

func (s *BlackholeSink) Id() *uuid.UUID {
	return s.id
}

func (s *BlackholeSink) Name() string {
	return s.name
}

func (s *BlackholeSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *BlackholeSink) Initialize(conf config.Sink) error {
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	return nil
}

func (s *BlackholeSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) error {
	return nil
}

func (s *BlackholeSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) error {
	return nil
}

func (s *BlackholeSink) Close() {
}
