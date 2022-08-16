// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

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
	deliveryRequired   bool
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

func (s *KafkaSink) Type() string {
	return KAFKA
}

func (s *KafkaSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *KafkaSink) Initialize(conf config.Sink) error {
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	ctx := context.Background()
	log.Debug().Msg("🟡 initializing kafka client")
	client, err := kgo.NewClient(
		kgo.SeedBrokers(conf.KafkaBrokers...),
	)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("could not create kafka sink client")
		return err
	}
	log.Debug().Msg("🟡 pinging kafka brokers")
	err = client.Ping(ctx)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("could not ping kafka sink brokers")
		return err
	}
	admClient := kadm.NewClient(client)
	log.Debug().Msg("🟡 verifying topics exist")
	topicDetails, err := admClient.DescribeTopicConfigs(ctx, conf.ValidTopic, conf.InvalidTopic)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not describe topic configs")
		return err
	}
	for _, d := range topicDetails {
		if d.Err != nil {
			log.Info().Msg("🟢 ensuring topic since it doesn't exist: " + d.Name)
			resp, err := admClient.CreateTopics(ctx, DEFAULT_PARTITIONS, DEFAULT_REPLICATION_FACTOR, nil, d.Name)
			if err != nil {
				log.Error().Err(err).Msg("🔴 could not create topic: " + d.Name)
			}
			log.Debug().Interface("response", resp).Msg("🟡 topic created: " + d.Name)
		}
	}
	s.client, s.validEventsTopic, s.invalidEventsTopic = client, conf.ValidTopic, conf.InvalidTopic
	return nil
}

func (s *KafkaSink) batchPublish(ctx context.Context, topic string, envelopes []envelope.Envelope) error {
	var wg sync.WaitGroup
	for _, e := range envelopes {
		payload, err := json.Marshal(e)
		if err != nil {
			return err
		}
		headers := []kgo.RecordHeader{
			{Key: envelope.INPUT_PROTOCOL, Value: []byte(e.EventMeta.Protocol)},
			{Key: envelope.VENDOR, Value: []byte(e.EventMeta.Vendor)},
			{Key: envelope.NAMESPACE, Value: []byte(e.EventMeta.Namespace)},
			{Key: envelope.VERSION, Value: []byte(e.EventMeta.Version)},
			{Key: envelope.FORMAT, Value: []byte(e.EventMeta.Format)},
		}
		record := &kgo.Record{
			Key:     []byte(e.EventMeta.Namespace),
			Topic:   topic,
			Value:   payload,
			Headers: headers,
		}
		wg.Add(1)
		var produceErr error // FIXME! Could probably do this similarly to pubsub sink and use a chan?
		s.client.Produce(ctx, record, func(r *kgo.Record, err error) {
			defer wg.Done()
			if err != nil {
				log.Error().Err(err).Msg("🔴 could not publish record")
				produceErr = err
			} else {
				offset := strconv.FormatInt(r.Offset, 10)
				partition := strconv.FormatInt(int64(r.Partition), 10)
				log.Trace().Msg("published event " + offset + " to topic " + topic + " partition " + partition)
			}
		})
		if produceErr != nil {
			return produceErr
		}
	}
	wg.Wait()
	return nil
}

func (s *KafkaSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validEventsTopic, envelopes)
	return err
}

func (s *KafkaSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.invalidEventsTopic, envelopes)
	return err
}

func (s *KafkaSink) Close() {
	log.Debug().Msg("🟡 closing kafka sink client")
	s.client.Close()
}
