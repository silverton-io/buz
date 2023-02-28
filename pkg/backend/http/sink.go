// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package http

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
)

type Sink struct {
	id               *uuid.UUID
	sinkType         string
	name             string
	deliveryRequired bool
	url              url.URL
	inputChan        chan []envelope.Envelope
	shutdown         chan int
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
	log.Debug().Msg("🟢 initializing " + s.sinkType + " sink")
	url, err := url.Parse(conf.Url)
	if err != nil {
		log.Debug().Err(err).Msg("🔴 " + conf.Url + " is not a valid url")
		return err
	}
	id := uuid.New()
	s.id, s.sinkType, s.name, s.deliveryRequired = &id, conf.Type, conf.Name, conf.DeliveryRequired
	s.url = *url
	s.inputChan = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
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
			case <-s.shutdown:
				return
			}
		}
	}(s)
	return nil
}

func (s *Sink) batchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	_, err := request.PostEnvelopes(s.url, envelopes)
	log.Error().Err(err).Interface("metadata", s.Metadata()).Msg("🔴 could not dequeue payloads")
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.inputChan <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	err := s.batchPublish(ctx, envelopes)
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Msg("🟢 shutting down " + s.sinkType + " sink") // no-op
	s.shutdown <- 1
	return nil
}
