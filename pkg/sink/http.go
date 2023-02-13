// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
)

type HttpSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	url              url.URL
}

func (s *HttpSink) Id() *uuid.UUID {
	return s.id
}

func (s *HttpSink) Name() string {
	return s.name
}

func (s *HttpSink) Type() string {
	return HTTP
}

func (s *HttpSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *HttpSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing http sink")
	url, err := url.Parse(conf.HttpUrl)
	if err != nil {
		log.Debug().Err(err).Msg("🟡 validUrl is not a valid url")
		return err
	}
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.url = *url
	return nil
}

func (s *HttpSink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	_, err := request.PostEnvelopes(s.url, envelopes)
	return err
}

func (s *HttpSink) Close() {
	log.Debug().Msg("🟡 closing http sink") // no-op
}
