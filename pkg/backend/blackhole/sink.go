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
	sinkType         string
	name             string
	deliveryRequired bool
	defaultOutput    string
	deadletterOutput string
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return backendutils.SinkMetadata{
		Id:               uuid.New(),
		Name:             s.name,
		SinkType:         s.sinkType,
		DeliveryRequired: s.deliveryRequired,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	id := uuid.New()
	s.id, s.sinkType, s.name = &id, conf.Type, conf.Name
	s.deliveryRequired = conf.DeliveryRequired
	return nil
}

func (s *Sink) StartWorker() error {
	// Blackhole. No worker necessary
	return nil
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	// This is a blackhole. It does nothing but dequeue
	ctx := context.Background()
	err := s.Dequeue(ctx, envelopes, "nothingness")
	if err != nil {
		log.Error().Err(err).Interface("metadata", s.Metadata()).Msg("could not dequeue")
	}
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	// This is a blackhole. It does nothing.
	return nil
}

func (s *Sink) Shutdown() error {
	log.Info().Msg("🟢 shutting down blackhole sink")
	return nil
}
