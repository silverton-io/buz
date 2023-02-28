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
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	id               *uuid.UUID
	sinkType         string
	name             string
	deliveryRequired bool
	fanout           bool
	outputFile       string
	inputChan        chan []envelope.Envelope
	shutdownChan     chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		SinkType:         s.sinkType,
		DeliveryRequired: s.deliveryRequired,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing file sink")
	id := uuid.New()
	s.id, s.name, s.sinkType = &id, conf.Name, conf.Type
	s.deliveryRequired, s.fanout = conf.DeliveryRequired, conf.Fanout
	s.inputChan = make(chan []envelope.Envelope, 10000)
	s.shutdownChan = make(chan int, 1)
	s.outputFile = "buz_events.json"
	s.StartWorker()
	return nil
}

func (s *Sink) StartWorker() error {
	go func(s *Sink) {
		for {
			select {
			case envelopes := <-s.inputChan:
				ctx := context.Background()
				s.Dequeue(ctx, envelopes)
			case <-s.shutdownChan:
				err := s.Shutdown()
				if err != nil {
					log.Error().Err(err).Interface("metadata", s.Metadata()).Msg("sink did not safely shut down")
				}
			}
		}
	}(s)
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

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.inputChan <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	err := s.batchPublish(ctx, s.outputFile, envelopes)
	return err
}

func (s *Sink) Shutdown() error {
	log.Info().Msg("🟢 shutting down file sink")
	s.shutdownChan <- 1
	return nil
}
