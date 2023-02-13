// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package pubnub

import (
	"context"
	"net/url"
	"path"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
)

const (
	PUBNUB_PUBLISH_URL   string = "ps.pndsn.com/publish"
	PUBNUB_HTTP_PROTOCOL string = "https"
)

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	validChannel     string
	invalidChannel   string
	pubKey           string
	subKey           string
	// store            int    // nolint: unused
	// callback         string // nolint: unused
}

func (s *Sink) Id() *uuid.UUID {
	return s.id
}

func (s *Sink) Name() string {
	return s.name
}

func (s *Sink) Type() string {
	return "pubnub"
}

func (s *Sink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing pubnub sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.validChannel, s.invalidChannel = constants.BUZ_VALID_EVENTS, constants.BUZ_INVALID_EVENTS
	s.pubKey, s.subKey = conf.PubnubPubKey, conf.PubnubSubKey
	return nil
}

func (s *Sink) buildPublishUrl(channel string) *url.URL {
	p := path.Join(PUBNUB_PUBLISH_URL, s.pubKey, s.subKey, "1", channel, "0")
	p = "https://" + p + "?uuid=" + s.id.String()
	u, err := url.Parse(p)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not parse publish url")
	}
	return u
}

func (s *Sink) batchPublish(ctx context.Context, channel string, envelopes []envelope.Envelope) error {
	u := s.buildPublishUrl(channel)
	_, err := request.PostEnvelopes(*u, envelopes)
	return err
}

func (s *Sink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validChannel, envelopes)
	return err
}

func (s *Sink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.invalidChannel, envelopes)
	return err
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	return nil
}

func (s *Sink) Close() {
	log.Debug().Msg("🟡 closing pubnub sink")
}
