package sink

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/net/context"
)

const (
	DEFAULT_PARTITIONS         int32 = 3
	DEFAULT_REPLICATION_FACTOR int16 = 3
)

type KafkaSink struct {
	id                 *uuid.UUID
	name               string
	client             *kgo.Client
	validEventsTopic   string
	invalidEventsTopic string
}

func (s *KafkaSink) Id() *uuid.UUID {
	return s.id
}

func (s *KafkaSink) Name() string {
	return s.name
}

func (s *KafkaSink) Initialize(conf config.Sink) error {
	ctx := context.Background()
	log.Debug().Msg("initializing kafka client")
	client, err := kgo.NewClient(
		kgo.SeedBrokers(conf.KafkaBrokers...),
	)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("could not create kafka sink client")
		return err
	}
	log.Debug().Msg("pinging kafka brokers")
	err = client.Ping(ctx)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("could not ping kafka sink brokers")
		return err
	}
	admClient := kadm.NewClient(client)
	log.Debug().Msg("verifying topics exist")
	topicDetails, err := admClient.DescribeTopicConfigs(ctx, conf.ValidEventTopic, conf.InvalidEventTopic)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("could not describe topic configs")
		return err
	}
	for _, d := range topicDetails {
		if d.Err != nil {
			log.Fatal().Stack().Err(d.Err).Msg("topic doesn't exist: " + d.Name)
		}
	}
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	s.client, s.validEventsTopic, s.invalidEventsTopic = client, conf.ValidEventTopic, conf.InvalidEventTopic
	return nil
}

func (s *KafkaSink) batchPublish(ctx context.Context, topic string, envelopes []envelope.Envelope) {
	var wg sync.WaitGroup
	for _, event := range envelopes {
		payload, _ := json.Marshal(event)
		headers := []kgo.RecordHeader{
			{Key: "vendor", Value: []byte(event.EventMetadata.Vendor)},
			{Key: "primaryNamespace", Value: []byte(event.EventMetadata.PrimaryNamespace)},
			{Key: "secondaryNamespace", Value: []byte(event.EventMetadata.SecondaryNamespace)},
			{Key: "tertiaryNamespace", Value: []byte(event.EventMetadata.TertiaryNamespace)},
			{Key: "name", Value: []byte(event.EventMetadata.Name)},
			{Key: "version", Value: []byte(event.EventMetadata.Version)},
		}
		record := &kgo.Record{
			Key:     []byte(event.EventMetadata.Path), // FIXME! Add configurable partition assignment
			Topic:   topic,
			Value:   payload,
			Headers: headers,
		}
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

func (s *KafkaSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.validEventsTopic, envelopes)
}

func (s *KafkaSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.invalidEventsTopic, envelopes)
}

func (s *KafkaSink) Close() {
	log.Debug().Msg("closing kafka sink client")
	s.client.Close()
}
