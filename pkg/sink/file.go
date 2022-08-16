// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package sink

import (
	"context"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type FileSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	validFile        string
	invalidFile      string
}

func (s *FileSink) Id() *uuid.UUID {
	return s.id
}

func (s *FileSink) Name() string {
	return s.name
}

func (s *FileSink) Type() string {
	return FILE
}

func (s *FileSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *FileSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing file sink")
	s.validFile = conf.ValidFile
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.invalidFile = conf.InvalidFile
	return nil
}

func (s *FileSink) batchPublish(ctx context.Context, filePath string, envelopes []envelope.Envelope) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open file")
		return err
	}
	for _, envelope := range envelopes {
		log.Debug().Msg("🟡 writing envelope to file " + filePath)
		b, err := json.Marshal(envelope)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not marshal envelope")
			return err
		}
		newline := []byte("\n")
		b = append(b, newline...)
		f.Write(b)
	}
	return nil
}

func (s *FileSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validFile, envelopes)
	return err
}

func (s *FileSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.invalidFile, envelopes)
	return err
}

func (s *FileSink) Close() {
	log.Debug().Msg("🟡 closing file sink")
}
