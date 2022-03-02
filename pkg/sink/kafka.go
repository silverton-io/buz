package sink

import (
	"encoding/json"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/net/context"
)

type KafkaSink struct {
	client             *kgo.Client
	validEventsTopic   string
	invalidEventsTopic string
}

func (s *KafkaSink) Initialize(conf config.Sink) {
	ctx := context.Background()
	log.Debug().Msg("initializing kafka client")
	client, err := kgo.NewClient(
		kgo.SeedBrokers(conf.Brokers...),
	)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not create kafka sink client")
	}
	log.Debug().Msg("pinging kafka brokers")
	err = client.Ping(ctx)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not ping kafka sink brokers")
	}
	admClient := kadm.NewClient(client)
	log.Debug().Msg("verifying topics exist")
	topicDetails, err := admClient.DescribeTopicConfigs(ctx, conf.ValidEventTopic, conf.InvalidEventTopic)
	for _, d := range topicDetails {
		if d.Err != nil {
			log.Fatal().Stack().Err(d.Err).Msg(d.Name + " topic doesn't exist")
		}
	}
	s.client, s.validEventsTopic, s.invalidEventsTopic = client, conf.ValidEventTopic, conf.InvalidEventTopic
}

func (s *KafkaSink) batchPublish(ctx context.Context, topic string, events []interface{}) {
	var wg sync.WaitGroup
	for _, event := range events {
		payload, _ := json.Marshal(event)
		record := &kgo.Record{Topic: topic, Value: payload}
		wg.Add(1)
		s.client.Produce(ctx, record, func(r *kgo.Record, err error) {
			defer wg.Done()
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not publish event to kafka")
			} else {
				offset := strconv.FormatInt(r.Offset, 10)
				partition := strconv.FormatInt(int64(r.Partition), 10)
				log.Debug().Msg("published event " + offset + " to topic " + topic + " partition " + partition)
			}
		})
	}
	wg.Wait()
}

func (s *KafkaSink) batchPublishValid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.validEventsTopic, events)
}

func (s *KafkaSink) batchPublishInvalid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.invalidEventsTopic, events)
}

func (s *KafkaSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta) {
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
	s.batchPublishValid(ctx, validEvents)
	s.batchPublishInvalid(ctx, invalidEvents)
	// Increment global metadata counters
	atomic.AddInt64(validCounter, int64(len(validEvents)))
	atomic.AddInt64(invalidCounter, int64(len(invalidEvents)))
}

func (s *KafkaSink) Close() {
	log.Debug().Msg("closing kafka sink client")
	s.client.Close()
}
