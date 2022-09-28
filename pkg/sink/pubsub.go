// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"encoding/json"
	"time"

	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"golang.org/x/net/context"
)

const INIT_TIMEOUT_SECONDS = 10

type PubsubSink struct {
	id                 *uuid.UUID
	name               string
	deliveryRequired   bool
	client             *pubsub.Client
	validEventsTopic   *pubsub.Topic
	invalidEventsTopic *pubsub.Topic
}

func (s *PubsubSink) Id() *uuid.UUID {
	return s.id
}

func (s *PubsubSink) Name() string {
	return s.name
}

func (s *PubsubSink) Type() string {
	return PUBSUB
}

func (s *PubsubSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *PubsubSink) Initialize(conf config.Sink) error {
	ctx, _ := context.WithTimeout(context.Background(), INIT_TIMEOUT_SECONDS*time.Second)
	client, err := pubsub.NewClient(ctx, conf.Project)
	if err != nil {
		log.Debug().Err(err).Msg("🟡 could not initialize pubsub sink")
		return err
	}
	validTopic := client.Topic(conf.ValidTopic)
	invalidTopic := client.Topic(conf.InvalidTopic)
	vTopicExists, err := validTopic.Exists(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("🟡 cannot check valid event topic existence")
		return err
	}
	if !vTopicExists {
		log.Debug().Err(err).Msg("🟡 valid event topic doesn't exist in project " + conf.Project)
		return err
	}
	invTopicExists, err := invalidTopic.Exists(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("🟡 cannot check invalid event topic existence")
		return err
	}
	if !invTopicExists {
		log.Debug().Err(err).Msg("🟡 invalid event topic doesn't exist in project " + conf.Project)
		return err
	}
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.client, s.validEventsTopic, s.invalidEventsTopic = client, validTopic, invalidTopic
	return nil
}

func (s *PubsubSink) batchPublish(ctx context.Context, topic *pubsub.Topic, envelopes []envelope.Envelope) error {
	var wg sync.WaitGroup
	for _, e := range envelopes {
		payload, _ := json.Marshal(e)
		msg := &pubsub.Message{
			Data: payload,
			Attributes: map[string]string{
				envelope.INPUT_PROTOCOL: e.EventMeta.Protocol,
				envelope.VENDOR:         e.EventMeta.Vendor,
				envelope.NAMESPACE:      e.EventMeta.Namespace,
				envelope.VERSION:        e.EventMeta.Version,
				envelope.FORMAT:         e.EventMeta.Format,
			},
		}
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

func (s *PubsubSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validEventsTopic, envelopes)
	return err
}

func (s *PubsubSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.invalidEventsTopic, envelopes)
	return err
}

func (s *PubsubSink) Close() {
	log.Debug().Msg("🟡 closing pubsub sink client")
	s.client.Close() // Technically does not need to be called since it's available for lifetime
}
