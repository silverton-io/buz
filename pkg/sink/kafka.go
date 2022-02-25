package sink

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/net/context"
)

type KafkaSink struct {
	client             *kgo.Client
	validEventsTopic   string
	invalidEventsTopic string
}

func (s *KafkaSink) Initialize(config config.Forwarder) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.Brokers...),
	)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not create kafka forwarder client")
	}
	s.client, s.validEventsTopic, s.invalidEventsTopic = client, config.ValidEventTopic, config.InvalidEventTopic
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

func (s *KafkaSink) BatchPublishValid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.validEventsTopic, events)
}

func (s *KafkaSink) BatchPublishInvalid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.invalidEventsTopic, events)
}

func (s *KafkaSink) Close() {
	log.Debug().Msg("closing kafka forwarder client")
	s.client.Close()
}
