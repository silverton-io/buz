// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package pubsub

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

const INIT_TIMEOUT_SECONDS = 10

type Sink struct {
	metadata backendutils.SinkMetadata
	client   *pubsub.Client
	input    chan []envelope.Envelope
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	ctx, cancel := context.WithTimeout(context.Background(), INIT_TIMEOUT_SECONDS*time.Second)
	defer cancel()
	client, err := pubsub.NewClient(ctx, conf.Project)
	if err != nil {
		log.Debug().Err(err).Msg("🟡 could not initialize pubsub sink")
		return err
	}
	s.client = client
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

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error {
	var wg sync.WaitGroup
	publisher := s.client.Publisher(output)
	defer publisher.Stop()
	for _, e := range envelopes {
		payload, _ := json.Marshal(e)
		msg := &pubsub.Message{
			Data: payload,
			Attributes: map[string]string{
				envelope.PROTOCOL:  e.Protocol,
				envelope.SCHEMA:    e.Schema,
				envelope.VENDOR:    e.Vendor,
				envelope.NAMESPACE: e.Namespace,
				envelope.VERSION:   e.Version,
				envelope.IS_VALID:  strconv.FormatBool(e.IsValid),
			},
		}
		result := publisher.Publish(ctx, msg)
		wg.Add(1)
		publishErr := make(chan error, 1)
		go func(res *pubsub.PublishResult, pErr chan error) {
			defer wg.Done()
			id, err := res.Get(ctx)
			if err != nil {
				pErr <- err

			} else {
				log.Trace().Msg("published event id " + id + " to topic " + publisher.ID())
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
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	s.shutdown <- 1
	s.client.Close() // Technically does not need to be called since it's available for lifetime
	return nil
}
