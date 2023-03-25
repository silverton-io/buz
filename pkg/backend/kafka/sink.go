// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package kafka

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"golang.org/x/net/context"
)

const (
	DEFAULT_PARTITIONS         int32 = 3
	DEFAULT_REPLICATION_FACTOR int16 = 1 // NOTE! Really not a good default.
)

type Sink struct {
	id                 *uuid.UUID
	sinkType           string
	name               string
	deliveryRequired   bool
	client             *kgo.Client
	defaultEventsTopic string
	input              chan []envelope.Envelope
	shutdown           chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		SinkType:         s.sinkType,
		DeliveryRequired: s.deliveryRequired,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	id := uuid.New()
	s.id, s.sinkType, s.name, s.deliveryRequired = &id, conf.Type, conf.Name, conf.DeliveryRequired
	ctx := context.Background()
	log.Debug().Msg("🟡 initializing kafka client")
	client, err := kgo.NewClient(
		kgo.SeedBrokers(conf.Brokers...),
	)
	s.client, s.defaultEventsTopic = client, constants.BUZ_EVENTS
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
	topicDetails, err := admClient.DescribeTopicConfigs(ctx, s.defaultEventsTopic)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not describe topic configs")
		return err
	}
	for _, d := range topicDetails {
		if d.Err != nil {
			log.Info().Msg("🟢 ensuring topic since it doesn't exist: " + d.Name)
			resp, err := admClient.CreateTopics(ctx, DEFAULT_PARTITIONS, DEFAULT_REPLICATION_FACTOR, nil, d.Name)
			if err != nil {
				log.Fatal().Err(err).Msg("🔴 could not create topic: " + d.Name)
			}
			log.Debug().Interface("response", resp).Msg("🟡 topic created: " + d.Name)
		}
	}
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	return nil
}

func (s *Sink) StartWorker() error {
	err := backendutils.StartSinkWorker(s.input, s.shutdown, s)
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.input <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	var wg sync.WaitGroup
	for _, e := range envelopes {
		payload, err := json.Marshal(e)
		if err != nil {
			return err
		}
		headers := []kgo.RecordHeader{
			{Key: envelope.PROTOCOL, Value: []byte(e.Protocol)},
			{Key: envelope.SCHEMA, Value: []byte(e.Schema)},
			{Key: envelope.VENDOR, Value: []byte(e.Vendor)},
			{Key: envelope.NAMESPACE, Value: []byte(e.Namespace)},
			{Key: envelope.VERSION, Value: []byte(e.Version)},
			{Key: envelope.IS_VALID, Value: []byte(strconv.FormatBool(e.IsValid))},
		}
		record := &kgo.Record{
			Key:     []byte(e.Namespace),
			Topic:   s.defaultEventsTopic,
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
				log.Trace().Msg("published event " + offset + " to topic " + s.defaultEventsTopic + " partition " + partition)
			}
		})
		if produceErr != nil {
			return produceErr
		}
	}
	wg.Wait()
	return nil
}

func (s *Sink) Shutdown() error {
	log.Debug().Msg("🟢 shutting down " + s.sinkType + " sink")
	s.shutdown <- 1
	s.client.Close()
	return nil
}
