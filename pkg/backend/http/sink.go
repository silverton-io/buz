// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package http

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
)

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	url              url.URL
}

func (s *Sink) Id() *uuid.UUID {
	return s.id
}

func (s *Sink) Name() string {
	return s.name
}

func (s *Sink) Type() string {
	return "http"
}

func (s *Sink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *Sink) Initialize(conf config.Sink) error {
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

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	_, err := request.PostEnvelopes(s.url, envelopes)
	return err
}

func (s *Sink) Close() {
	log.Debug().Msg("🟡 closing http sink") // no-op
}
