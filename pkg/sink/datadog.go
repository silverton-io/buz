// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type DatadogSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
}

func (s *DatadogSink) Id() *uuid.UUID {
	return s.id
}

func (s *DatadogSink) Name() string {
	return s.name
}

func (s *DatadogSink) Type() string {
	return DATADOG
}

func (s *DatadogSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *DatadogSink) Initialize(conf config.Sink) error {
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	return nil
}

func (s *DatadogSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) error {
	return nil
}

func (s *DatadogSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) error {
	return nil
}

func (s *DatadogSink) Close() {
}
