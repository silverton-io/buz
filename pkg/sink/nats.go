package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type NatsSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	validSubject     string
	invalidSubject   string
}

func (s *NatsSink) Id() *uuid.UUID {
	return s.id
}

func (s *NatsSink) Name() string {
	return s.name
}

func (s *NatsSink) Type() string {
	return NATS
}

func (s *NatsSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *NatsSink) Initialize() error {
	// FIXME!
	return nil
}

func (s *NatsSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	// FIXME!
	return nil
}

func (s *NatsSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	// FIXME!
	return nil
}

func (s *NatsSink) Close() {
	log.Debug().Msg("closing nats sink")
}
