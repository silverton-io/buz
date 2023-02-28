// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package kinesis

import (
	"context"
	"encoding/json"
	"sync"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	client           *kinesis.Client
	validStream      string
	invalidStream    string
}

func (s *Sink) Id() *uuid.UUID {
	return s.id
}

func (s *Sink) Name() string {
	return s.name
}

func (s *Sink) Type() string {
	return "kinesis"
}

func (s *Sink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *Sink) Initialize(conf config.Sink) error {
	ctx := context.Background()
	cfg, err := awsconf.LoadDefaultConfig(ctx)
	client := kinesis.NewFromConfig(cfg)
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.client, s.validStream, s.invalidStream = client, constants.BUZ_VALID_EVENTS, constants.BUZ_INVALID_EVENTS
	return err
}

func (s *Sink) batchPublish(ctx context.Context, stream string, envelopes []envelope.Envelope) error {
	var wg sync.WaitGroup
	for _, event := range envelopes {
		partitionKey := "blah" // FIXME!
		payload, _ := json.Marshal(event)
		input := &kinesis.PutRecordInput{
			Data:         payload,
			PartitionKey: &partitionKey,
			StreamName:   &stream,
		}
		wg.Add(1)
		pubErr := make(chan error, 1)
		go func(pErr chan error) {
			output, err := s.client.PutRecord(ctx, input) // Will want to use `PutRecordBatch`
			defer wg.Done()
			if err != nil {
				log.Error().Err(err).Msg("🔴 could not publish event to kinesis")
				pErr <- err
			} else {
				log.Debug().Msgf("🟡 published event " + *output.SequenceNumber + " to stream " + stream)
				pErr <- nil
			}
		}(pubErr)
		err := <-pubErr
		if err != nil {
			return err
		}
	}
	wg.Wait()
	return nil
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validStream, envelopes) // FIXME -> shard by configured strategy
	return err
}

func (s *Sink) Close() {
	log.Debug().Msg("🟡 closing kinesis sink client")

}