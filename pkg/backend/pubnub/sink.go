// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package pubnub

import (
	"context"
	"net/url"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
)

const (
	PUBNUB_PUBLISH_URL   string = "ps.pndsn.com/publish"
	PUBNUB_HTTP_PROTOCOL string = "https"
)

type Sink struct {
	metadata backendutils.SinkMetadata
	pubKey   string
	subKey   string
	input    chan []envelope.Envelope
	shutdown chan int
	// store            int    // nolint: unused
	// callback         string // nolint: unused
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return s.metadata
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing pubnub sink")
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	s.pubKey, s.subKey = conf.PubnubPubKey, conf.PubnubSubKey
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

func (s *Sink) buildPublishUrl(channel string) *url.URL {
	p := path.Join(PUBNUB_PUBLISH_URL, s.pubKey, s.subKey, "1", channel, "0")
	p = "https://" + p + "?uuid=" + s.metadata.Id.String()
	u, err := url.Parse(p)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not parse publish url")
	}
	return u
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error {
	u := s.buildPublishUrl(output)
	_, err := request.PostEnvelopes(*u, envelopes)
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Interface("metadata", s.metadata).Msg("🟢 shutting down sink")
	s.shutdown <- 1
	return nil
}
