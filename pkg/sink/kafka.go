package sink

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/net/context"
)

const (
	DEFAULT_PARTITIONS         int32 = 3
	DEFAULT_REPLICATION_FACTOR int16 = 3
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
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not describe topic configs")
	}
	for _, d := range topicDetails {
		if d.Err != nil {
			log.Fatal().Stack().Err(d.Err).Msg("topic doesn't exist: " + d.Name)
		}
	}
	s.client, s.validEventsTopic, s.invalidEventsTopic = client, conf.ValidEventTopic, conf.InvalidEventTopic
}

func (s *KafkaSink) batchPublish(ctx context.Context, topic string, envelopes []envelope.Envelope) {
	var wg sync.WaitGroup
	for _, event := range envelopes {
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

func (s *KafkaSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.validEventsTopic, envelopes)
}

func (s *KafkaSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.invalidEventsTopic, envelopes)
}

func (s *KafkaSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope, meta *tele.Meta) {
	// Publish
	go s.BatchPublishValid(ctx, validEnvelopes)
	go s.BatchPublishInvalid(ctx, invalidEnvelopes)
	// Increment stats counters
	incrementStats(inputType, len(validEnvelopes), len(invalidEnvelopes), meta)
}

func (s *KafkaSink) Close() {
	log.Debug().Msg("closing kafka sink client")
	s.client.Close()
}
