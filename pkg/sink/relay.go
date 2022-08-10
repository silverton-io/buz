// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package sink

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/request"
)

type RelaySink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	relayUrl         url.URL
}

func (s *RelaySink) Id() *uuid.UUID {
	return s.id
}

func (s *RelaySink) Name() string {
	return s.name
}

func (s *RelaySink) Type() string {
	return RELAY
}

func (s *RelaySink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *RelaySink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing http sink")
	u, err := url.Parse(conf.RelayUrl)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("relayUrl is not a valid url")
		return err
	}
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.relayUrl = *u
	return err
}

func (s *RelaySink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) error {
	_, err := request.PostEnvelopes(s.relayUrl, validEnvelopes)
	return err
}

func (s *RelaySink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) error {
	_, err := request.PostEnvelopes(s.relayUrl, invalidEnvelopes)
	return err
}

func (s *RelaySink) Close() {
	log.Debug().Msg("closing relay sink") // no-op
}
