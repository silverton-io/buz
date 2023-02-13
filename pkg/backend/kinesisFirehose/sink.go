// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package kinesisFirehose

import (
	"context"
	"encoding/json"
	"sync"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
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
	client           *firehose.Client
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
	return "firehose"
}

func (s *Sink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *Sink) Initialize(conf config.Sink) error {
	ctx := context.Background()
	cfg, err := awsconf.LoadDefaultConfig(ctx)
	client := firehose.NewFromConfig(cfg)
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.client, s.validStream, s.invalidStream = client, constants.BUZ_VALID_EVENTS, constants.BUZ_INVALID_EVENTS
	return err
}

func (s *Sink) batchPublish(ctx context.Context, stream string, envelopes []envelope.Envelope) error {
	var wg sync.WaitGroup
	var records []types.Record
	for _, event := range envelopes {
		payload, _ := json.Marshal(event)
		newline := []byte("\n")
		payload = append(payload, newline...)
		record := types.Record{
			Data: payload,
		}
		records = append(records, record)
	}
	input := &firehose.PutRecordBatchInput{
		DeliveryStreamName: &stream,
		Records:            records,
	}
	wg.Add(1)
	pubErr := make(chan error, 1)
	go func(pErr chan error) {
		_, err := s.client.PutRecordBatch(ctx, input)
		defer wg.Done()
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not publish event to kinesis firehose")
			pubErr <- err
		} else {
			log.Debug().Msgf("🟡 published event batch to stream " + stream)
			pubErr <- nil
		}
	}(pubErr)
	err := <-pubErr
	if err != nil {
		return err
	}
	wg.Wait()
	return nil
}

func (s *Sink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validStream, envelopes)
	return err
}

func (s *Sink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.invalidStream, envelopes)
	return err
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	// err := s.batchPublish(ctx, s.validStream, envelopes)
	return nil
}

func (s *Sink) Close() {
	log.Debug().Msg("🟡 closing kinesis firehose sink client")
}
