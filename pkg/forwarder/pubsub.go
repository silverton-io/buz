package forwarder

import (
	"encoding/json"

	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"golang.org/x/net/context"
)

type PubsubForwarder struct {
	Client             *pubsub.Client
	validEventsTopic   *pubsub.Topic
	invalidEventsTopic *pubsub.Topic
}

func (f *PubsubForwarder) Initialize(config config.Forwarder) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, config.Project)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not initialize forwarder")
	}
	validTopic := client.Topic(config.ValidEventTopic)
	invalidTopic := client.Topic(config.InvalidEventTopic)
	f.Client, f.validEventsTopic, f.invalidEventsTopic = client, validTopic, invalidTopic
}

func (f *PubsubForwarder) publishEvent(ctx context.Context, topic *pubsub.Topic, event interface{}) {
	payload, _ := json.Marshal(event)
	msg := &pubsub.Message{
		Data: payload,
	}
	result := topic.Publish(ctx, msg)
	id, err := result.Get(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not publish event")
	} else {
		log.Debug().Msgf("published event id " + id + " to topic " + topic.ID())
	}
}

func (f *PubsubForwarder) batchPublishEvents(ctx context.Context, topic *pubsub.Topic, events []interface{}) {
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
				log.Error().Stack().Err(err).Msg("could not publish event")
			} else {
				log.Debug().Msgf("published event id " + id + " to topic " + topic.ID())
			}
		}(result)
	}
	wg.Wait()
}

func (f *PubsubForwarder) PublishValidEvent(ctx context.Context, event interface{}) {
	f.publishEvent(ctx, f.validEventsTopic, event)
}

func (f *PubsubForwarder) PublishInvalidEvent(ctx context.Context, event interface{}) {
	f.publishEvent(ctx, f.invalidEventsTopic, event)
}

func (f *PubsubForwarder) PublishValidEvents(ctx context.Context, events []interface{}) {

	f.batchPublishEvents(ctx, f.validEventsTopic, events)
}

func (f *PubsubForwarder) PublishInvalidEvents(ctx context.Context, events []interface{}) {
	f.batchPublishEvents(ctx, f.invalidEventsTopic, events)
}
