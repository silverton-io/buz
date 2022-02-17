package forwarder

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/net/context"
)

type KafkaForwarder struct {
	client             *kgo.Client
	validEventsTopic   string
	invalidEventsTopic string
}

func (f *KafkaForwarder) Initialize(config config.Forwarder) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(config.Brokers...),
	)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not create kafka forwarder client")
	}
	f.client, f.validEventsTopic, f.invalidEventsTopic = client, config.ValidEventTopic, config.InvalidEventTopic
}

func (f *KafkaForwarder) batchPublish(ctx context.Context, topic string, events []interface{}) {
	var wg sync.WaitGroup
	for _, event := range events {
		payload, _ := json.Marshal(event)
		record := &kgo.Record{Topic: topic, Value: payload}
		wg.Add(1)
		f.client.Produce(ctx, record, func(r *kgo.Record, err error) {
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

func (f *KafkaForwarder) BatchPublishValid(ctx context.Context, events []interface{}) {
	f.batchPublish(ctx, f.validEventsTopic, events)
}

func (f *KafkaForwarder) BatchPublishInvalid(ctx context.Context, events []interface{}) {
	f.batchPublish(ctx, f.invalidEventsTopic, events)
}

func (f *KafkaForwarder) Close() {
	log.Debug().Msg("closing kafka forwarder client")
	f.client.Close()
}
