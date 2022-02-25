package sink

import (
	"encoding/json"
	"sync/atomic"

	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"golang.org/x/net/context"
)

type PubsubSink struct {
	client             *pubsub.Client
	validEventsTopic   *pubsub.Topic
	invalidEventsTopic *pubsub.Topic
}

func (s *PubsubSink) Initialize(config config.Forwarder) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, config.Project)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not initialize forwarder")
	}
	validTopic := client.Topic(config.ValidEventTopic)
	invalidTopic := client.Topic(config.InvalidEventTopic)
	s.client, s.validEventsTopic, s.invalidEventsTopic = client, validTopic, invalidTopic
}

func (s *PubsubSink) publish(ctx context.Context, topic *pubsub.Topic, event interface{}) {
	var events []interface{}
	events = append(events, event)
	s.batchPublish(ctx, topic, events)
}

func (s *PubsubSink) batchPublish(ctx context.Context, topic *pubsub.Topic, events []interface{}) {
	var wg sync.WaitGroup
	for _, event := range events {
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

func (s *PubsubSink) PublishValid(ctx context.Context, event interface{}) {
	s.publish(ctx, s.validEventsTopic, event)
}

func (s *PubsubSink) PublishInvalid(ctx context.Context, event interface{}) {
	s.publish(ctx, s.invalidEventsTopic, event)
}

func (s *PubsubSink) BatchPublishValid(ctx context.Context, events []interface{}) {

	s.batchPublish(ctx, s.validEventsTopic, events)
}

func (s *PubsubSink) BatchPublishInvalid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.invalidEventsTopic, events)
}

func (s *PubsubSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta) {
	var validCounter *int64
	var invalidCounter *int64
	switch inputType {
	case input.GENERIC_INPUT:
		validCounter = &meta.ValidGenericEventsProcessed
		invalidCounter = &meta.InvalidGenericEventsProcessed
	case input.CLOUDEVENTS_INPUT:
		validCounter = &meta.ValidCloudEventsProcessed
		invalidCounter = &meta.InvalidCloudEventsProcessed
	default:
		validCounter = &meta.ValidSnowplowEventsProcessed
		invalidCounter = &meta.InvalidSnowplowEventsProcessed
	}
	// Publish
	s.BatchPublishValid(ctx, validEvents)
	s.BatchPublishInvalid(ctx, invalidEvents)
	// Increment global metadata counters
	atomic.AddInt64(validCounter, int64(len(validEvents)))
	atomic.AddInt64(invalidCounter, int64(len(invalidEvents)))
}

func (s *PubsubSink) Close() {
	log.Debug().Msg("closing pubsub forwarder client")
	s.client.Close() // Technically does not need to be called since it's available for lifetime
}
