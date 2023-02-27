// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package blackhole

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	fanout           bool
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	sinkType := "blackhole"
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		Type:             sinkType,
		DeliveryRequired: s.deliveryRequired,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	s.deliveryRequired, s.fanout = conf.DeliveryRequired, conf.Fanout
	return nil
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	// This is a blackhole. It does nothing.
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	// This is a blackhole. It does nothing.
	return nil
}

func (s *Sink) Shutdown() error {
	log.Info().Msg("🟢 shutting down blackhole sink")
	return nil
}
