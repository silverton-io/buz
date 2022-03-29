package sink

import (
	"encoding/json"
	"time"

	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"golang.org/x/net/context"
)

const INIT_TIMEOUT_SECONDS = 10

type PubsubSink struct {
	id                 *uuid.UUID
	name               string
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

func (s *PubsubSink) Initialize(conf config.Sink) error {
	ctx, _ := context.WithTimeout(context.Background(), INIT_TIMEOUT_SECONDS*time.Second)
	client, err := pubsub.NewClient(ctx, conf.Project)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("could not initialize pubsub sink")
		return err
	}
	validTopic := client.Topic(conf.ValidEventTopic)
	invalidTopic := client.Topic(conf.InvalidEventTopic)
	vTopicExists, err := validTopic.Exists(ctx)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("cannot check valid event topic existence")
		return err
	}
	if !vTopicExists {
		log.Debug().Stack().Err(err).Msg("valid event topic doesn't exist in project " + conf.Project)
		return err
	}
	invTopicExists, err := invalidTopic.Exists(ctx)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("cannot check invalid event topic existence")
		return err
	}
	if !invTopicExists {
		log.Debug().Stack().Err(err).Msg("invalid event topic doesn't exist in project " + conf.Project)
		return err
	}
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	s.client, s.validEventsTopic, s.invalidEventsTopic = client, validTopic, invalidTopic
	return nil
}

func (s *PubsubSink) batchPublish(ctx context.Context, topic *pubsub.Topic, envelopes []envelope.Envelope) {
	var wg sync.WaitGroup
	for _, event := range envelopes {
		payload, _ := json.Marshal(event)
		msg := &pubsub.Message{
			Data: payload,
		}
		result := topic.Publish(ctx, msg)
		wg.Add(1)
		go func(res *pubsub.PublishResult) {
			defer wg.Done()
			id, err := res.Get(ctx)
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not publish event to pubsub")
			} else {
				log.Debug().Msgf("published event id " + id + " to topic " + topic.ID())
			}
		}(result)
	}
	wg.Wait()
}

func (s *PubsubSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.validEventsTopic, envelopes)
}

func (s *PubsubSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.invalidEventsTopic, envelopes)
}

func (s *PubsubSink) Close() {
	log.Debug().Msg("closing pubsub sink client")
	s.client.Close() // Technically does not need to be called since it's available for lifetime
}
