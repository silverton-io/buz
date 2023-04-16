// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package splunk

import (
	"context"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	url      *url.URL
	apiKey   string
	input    chan []envelope.Envelope
	shutdown chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	url, err := url.Parse(conf.Url)
	if err != nil {
		log.Fatal().Err(err).Interface("metadata", s.Metadata()).Msg(conf.Url + " is not a valid url")
	}
	s.url = url
	s.apiKey = conf.ApiKey
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
	splunkHeader := http.Header{
		"Authorization": {"Splunk " + s.apiKey},
	}
	resp, err := request.PostEnvelopes(*s.url, envelopes, splunkHeader)
	if err != nil {
		log.Error().Interface("response", resp).Err(err).Msg("could not post envelopes")
	}
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	s.shutdown <- 1
	return nil
}
