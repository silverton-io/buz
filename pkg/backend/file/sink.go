// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package file

import (
	"context"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	fanout           bool
	validFile        string
	invalidFile      string
	inputChan        chan []envelope.Envelope
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	sinkType := "file"
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		Type:             sinkType,
		DeliveryRequired: s.deliveryRequired,
		Fanout:           false,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing file sink")
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	s.deliveryRequired, s.fanout = conf.DeliveryRequired, conf.Fanout
	s.inputChan = make(chan []envelope.Envelope, 10000)
	s.validFile = constants.BUZ_VALID_EVENTS + ".json"
	s.invalidFile = constants.BUZ_INVALID_EVENTS + ".json"
	return nil
}

func (s *Sink) batchPublish(ctx context.Context, filePath string, envelopes []envelope.Envelope) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open file")
		return err
	}
	defer f.Close() // nolint
	for _, envelope := range envelopes {
		log.Debug().Msg("🟡 writing envelope to file " + filePath)
		b, err := json.Marshal(envelope)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not marshal envelope")
			return err
		}
		newline := []byte("\n")
		b = append(b, newline...)
		if _, err := f.Write(b); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validFile, envelopes) // FIXME -> shard this by configured strategy
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) {
	s.inputChan <- envelopes
}

func (s *Sink) Shutdown() error {
	log.Info().Msg("🟡 shutting down file sink")
	return nil
}
