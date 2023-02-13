// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package blackhole

import (
	"context"

	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
}

func (s *Sink) Id() *uuid.UUID {
	return s.id
}

func (s *Sink) Name() string {
	return s.name
}

func (s *Sink) Type() string {
	return "blackhole"
}

func (s *Sink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *Sink) Initialize(conf config.Sink) error {
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	return nil
}

func (s *Sink) BatchPublish(ctx context.Context, validEnvelopes []envelope.Envelope) error {
	return nil
}

func (s *Sink) Close() {
}
