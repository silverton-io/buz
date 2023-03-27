// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package http

import (
	"context"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	input    chan []envelope.Envelope
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	_, err := url.Parse(conf.DefaultOutput)
	if err != nil {
		log.Error().Err(err).Msg("🔴 " + conf.DefaultOutput + " is not a valid url")
		return err
	}
	_, err = url.Parse(conf.DeadletterOutput)
	if err != nil {
		log.Error().Err(err).Msg("🔴 " + conf.DeadletterOutput + " is not a valid url")
		return err
	}
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
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
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	url, err := url.Parse(output)
	if err != nil {
		log.Error().Err(err).Msg("🔴 " + output + " is not a valid url")
		return err
	}
	_, err = request.PostEnvelopes(*url, envelopes)
	if err != nil {
		log.Error().Err(err).Interface("metadata", s.Metadata()).Msg("🔴 could not dequeue payloads")
	}
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink") // no-op
	s.shutdown <- 1
	return nil
}
