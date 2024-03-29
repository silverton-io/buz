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
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	client   *firehose.Client
	input    chan []envelope.Envelope
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	ctx := context.Background()
	cfg, err := awsconf.LoadDefaultConfig(ctx)
	client := firehose.NewFromConfig(cfg)
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	s.client = client
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	return err
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

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
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
		DeliveryStreamName: &output,
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
			log.Debug().Msgf("🟡 published event batch to stream " + *input.DeliveryStreamName)
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

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	s.shutdown <- 1
	return nil
}
