// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package file

import (
	"context"
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	input    chan []envelope.Envelope
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing file sink")
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	return nil
}

func (s *Sink) StartWorker() error {
	err := backendutils.StartSinkWorker(s.input, s.shutdown, s)
	return err
}

func (s *Sink) batchPublish(ctx context.Context, filePath string, envelopes []envelope.Envelope) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open file")
		return err
	}
	defer f.Close() // nolint
	for _, envelope := range envelopes {
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

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.input <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	err := s.batchPublish(ctx, output, envelopes)
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	s.shutdown <- 1
	return nil
}
