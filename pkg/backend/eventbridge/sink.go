// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package eventbridge

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	client   *eventbridge.EventBridge
	input    chan []envelope.Envelope
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	es := session.Must(session.NewSession())
	svc := eventbridge.New(es)
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	s.client = svc
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

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error {
	var entries []*eventbridge.PutEventsRequestEntry
	for _, e := range envelopes {
		byteString, err := e.AsByte()
		if err != nil {
			log.Error().Err(err).Msg("could not cast envelope to bytes")
		}
		entry := eventbridge.PutEventsRequestEntry{
			EventBusName: &output,
			Time:         &e.Timestamp,
			Source:       aws.String(constants.BUZ),
			DetailType:   &e.Schema,
			Detail:       aws.String(string(byteString)),
		}
		entries = append(entries, &entry)
	}
	input := eventbridge.PutEventsInput{
		Entries: entries,
	}
	result, err := s.client.PutEvents(&input)
	if err != nil {
		log.Error().Err(err).Interface("result", result).Msg("could not dequeue")
		return err
	}
	return nil
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	s.shutdown <- 1
	return nil
}
