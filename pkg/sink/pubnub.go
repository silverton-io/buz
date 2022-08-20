// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"context"
	"net/url"
	"path"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
)

const (
	PUBNUB_PUBLISH_URL   string = "ps.pndsn.com/publish"
	PUBNUB_HTTP_PROTOCOL string = "https"
)

type PubnubSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	validChannel     string
	invalidChannel   string
	pubKey           string
	subKey           string
	store            int
	callback         string
}

func (s *PubnubSink) Id() *uuid.UUID {
	return s.id
}

func (s *PubnubSink) Name() string {
	return s.name
}

func (s *PubnubSink) Type() string {
	return PUBNUB
}

func (s *PubnubSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *PubnubSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing pubnub sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.validChannel, s.invalidChannel = conf.ValidChannel, conf.InvalidChannel
	s.pubKey, s.subKey = conf.PubnubPubKey, conf.PubnubSubKey
	return nil
}

func (s *PubnubSink) buildPublishUrl(channel string) *url.URL {
	p := path.Join(PUBNUB_PUBLISH_URL, s.pubKey, s.subKey, "1", channel, "0")
	p = "https://" + p + "?uuid=" + s.id.String()
	u, err := url.Parse(p)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not parse publish url")
	}
	return u
}

func (s *PubnubSink) batchPublish(ctx context.Context, channel string, envelopes []envelope.Envelope) error {
	u := s.buildPublishUrl(channel)
	_, err := request.PostEnvelopes(*u, envelopes)
	return err
}

func (s *PubnubSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validChannel, envelopes)
	return err
}

func (s *PubnubSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.invalidChannel, envelopes)
	return err
}

func (s *PubnubSink) Close() {
	log.Debug().Msg("🟡 closing pubnub sink")
}
