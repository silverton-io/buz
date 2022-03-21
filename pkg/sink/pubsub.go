package sink

import (
	"encoding/json"
	"time"

	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/tele"
	"golang.org/x/net/context"
)

type PubsubSink struct {
	client             *pubsub.Client
	validEventsTopic   *pubsub.Topic
	invalidEventsTopic *pubsub.Topic
}

const INIT_TIMEOUT_SECONDS = 10

func (s *PubsubSink) Initialize(conf config.Sink) {
	ctx, _ := context.WithTimeout(context.Background(), INIT_TIMEOUT_SECONDS*time.Second)
	client, err := pubsub.NewClient(ctx, conf.Project)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not initialize pubsub sink")
	}
	validTopic := client.Topic(conf.ValidEventTopic)
	invalidTopic := client.Topic(conf.InvalidEventTopic)
	vTopicExists, err := validTopic.Exists(ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("cannot check valid event topic existence")
	}
	if !vTopicExists {
		log.Fatal().Stack().Err(err).Msg("valid event topic doesn't exist in project " + conf.Project)
	}
	invTopicExists, err := invalidTopic.Exists(ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("cannot check invalid event topic existence")
	}
	if !invTopicExists {
		log.Fatal().Stack().Err(err).Msg("invalid event topic doesn't exist in project " + conf.Project)
	}
	s.client, s.validEventsTopic, s.invalidEventsTopic = client, validTopic, invalidTopic
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

func (s *PubsubSink) BatchPublishValidAndInvalid(ctx context.Context, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope, meta *tele.Meta) {
	// Publish
	go s.BatchPublishValid(ctx, validEnvelopes)
	go s.BatchPublishInvalid(ctx, invalidEnvelopes)
	// Increment stats counters
}

func (s *PubsubSink) Close() {
	log.Debug().Msg("closing pubsub sink client")
	s.client.Close() // Technically does not need to be called since it's available for lifetime
}
