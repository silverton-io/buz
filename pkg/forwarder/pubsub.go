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
	client             *pubsub.Client
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
	f.client, f.validEventsTopic, f.invalidEventsTopic = client, validTopic, invalidTopic
}

func (f *PubsubForwarder) publish(ctx context.Context, topic *pubsub.Topic, event interface{}) {
	var events []interface{}
	events = append(events, event)
	f.batchPublish(ctx, topic, events)
}

func (f *PubsubForwarder) batchPublish(ctx context.Context, topic *pubsub.Topic, events []interface{}) {
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

func (f *PubsubForwarder) PublishValid(ctx context.Context, event interface{}) {
	f.publish(ctx, f.validEventsTopic, event)
}

func (f *PubsubForwarder) PublishInvalid(ctx context.Context, event interface{}) {
	f.publish(ctx, f.invalidEventsTopic, event)
}

func (f *PubsubForwarder) BatchPublishValid(ctx context.Context, events []interface{}) {

	f.batchPublish(ctx, f.validEventsTopic, events)
}

func (f *PubsubForwarder) BatchPublishInvalid(ctx context.Context, events []interface{}) {
	f.batchPublish(ctx, f.invalidEventsTopic, events)
}

func (f *PubsubForwarder) Close() {
	log.Debug().Msg("closing pubsub forwarder client")
	f.client.Close() // Technically does not need to be called since it's available for lifetime
}
