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
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"golang.org/x/net/context"
)

const INIT_TIMEOUT_SECONDS = 10

// See whether or not a topic exists.
func checkTopicExistence(topic *pubsub.Topic) bool {
	ctx := context.Background()
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("🟡 cannot check valid event topic existence")
		return false
	}
	if !exists {
		return false
	}
	return true
}

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
	ctx, _ := context.WithTimeout(context.Background(), INIT_TIMEOUT_SECONDS*time.Second)
	client, err := pubsub.NewClient(ctx, conf.Project)
	if err != nil {
		log.Debug().Err(err).Msg("🟡 could not initialize pubsub sink")
		return err
	}
	defaultTopic := client.Topic(s.metadata.DefaultOutput)
	deadletterTopic := client.Topic(s.metadata.DeadletterOutput)
	for _, i := range []pubsub.Topic{*defaultTopic, *deadletterTopic} {
		checkTopicExistence(&i)
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
		topic := s.client.Topic(output)
		result := topic.Publish(ctx, msg)
		wg.Add(1)
		publishErr := make(chan error, 1)
		go func(res *pubsub.PublishResult, pErr chan error) {
			defer wg.Done()
			id, err := res.Get(ctx)
			if err != nil {
				pErr <- err

			} else {
				log.Trace().Msgf("published event id " + id + " to topic " + topic.ID())
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
