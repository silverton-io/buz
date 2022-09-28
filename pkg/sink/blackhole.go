// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type BlackholeSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
}

func (s *BlackholeSink) Id() *uuid.UUID {
	return s.id
}

func (s *BlackholeSink) Name() string {
	return s.name
}

func (s *BlackholeSink) Type() string {
	return BLACKHOLE
}

func (s *BlackholeSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *BlackholeSink) Initialize(conf config.Sink) error {
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	return nil
}

func (s *BlackholeSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) error {
	return nil
}

func (s *BlackholeSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) error {
	return nil
}

func (s *BlackholeSink) Close() {
}
