// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package pubsub

import (
	"encoding/json"
	"strconv"
	"time"

	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
	"golang.org/x/net/context"
)

const INIT_TIMEOUT_SECONDS = 10

type Sink struct {
	id                 *uuid.UUID
	sinkType           string
	name               string
	deliveryRequired   bool
	client             *pubsub.Client
	defaultEventsTopic *pubsub.Topic
	input              chan []envelope.Envelope
	shutdown           chan int
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
	ctx, _ := context.WithTimeout(context.Background(), INIT_TIMEOUT_SECONDS*time.Second)
	client, err := pubsub.NewClient(ctx, conf.Project)
	if err != nil {
		log.Debug().Err(err).Msg("🟡 could not initialize pubsub sink")
		return err
	}
	defaultTopic := client.Topic(constants.BUZ_EVENTS)
	dTopicExists, err := defaultTopic.Exists(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("🟡 cannot check valid event topic existence")
		return err
	}
	if !dTopicExists {
		log.Debug().Err(err).Msg("🟡 valid event topic doesn't exist in project " + conf.Project)
		return err
	}
	id := uuid.New()
	s.id, s.sinkType, s.name, s.deliveryRequired = &id, conf.Type, conf.Name, conf.DeliveryRequired
	s.client, s.defaultEventsTopic = client, defaultTopic
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	return nil
}

func (s *Sink) StartWorker() error {
	err := backendutils.StartSinkWorker(s.input, s.shutdown, s)
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.input <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	var wg sync.WaitGroup
	for _, e := range envelopes {
		payload, _ := json.Marshal(e)
		msg := &pubsub.Message{
			Data: payload,
			Attributes: map[string]string{
				envelope.PROTOCOL:  e.Protocol,
				envelope.NAMESPACE: e.Namespace,
				envelope.VERSION:   e.Version,
				envelope.FORMAT:    e.Format,
				envelope.SCHEMA:    e.Schema,
				envelope.IS_VALID:  strconv.FormatBool(e.IsValid),
			},
		}
		result := s.defaultEventsTopic.Publish(ctx, msg)
		wg.Add(1)
		publishErr := make(chan error, 1)
		go func(res *pubsub.PublishResult, pErr chan error) {
			defer wg.Done()
			id, err := res.Get(ctx)
			if err != nil {
				pErr <- err

			} else {
				log.Trace().Msgf("published event id " + id + " to topic " + s.defaultEventsTopic.ID())
				pErr <- nil
			}
		}(result, publishErr)
		err := <-publishErr
		if err != nil {
			return err
		}
	}
	wg.Wait()
	return nil
}

func (s *Sink) Shutdown() error {
	log.Debug().Msg("🟢 shutting down " + s.sinkType + " sink")
	s.shutdown <- 1
	s.client.Close() // Technically does not need to be called since it's available for lifetime
	return nil
}
