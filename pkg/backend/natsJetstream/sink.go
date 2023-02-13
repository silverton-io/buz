// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package natsJetstream

import (
	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	conn             *nats.Conn
	jetstream        nats.JetStreamContext
	validSubject     string
	invalidSubject   string
	// FIXME! Add .creds/token/tls cert/nkey auth
}

func (s *Sink) Id() *uuid.UUID {
	return s.id
}

func (s *Sink) Name() string {
	return s.name
}

func (s *Sink) Type() string {
	return "jetstream"
}

func (s *Sink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing nats jetstream sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	conn, err := nats.Connect(conf.NatsHost, nats.UserInfo(conf.NatsUser, conf.NatsPass))
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open nats connection")
		return err
	}
	js, err := conn.JetStream()

	if err != nil {
		log.Error().Err(err).Msg("🔴 could not use jetstream context")
		return err
	}

	s.validSubject, s.invalidSubject = constants.BUZ_VALID_EVENTS, constants.BUZ_INVALID_EVENTS
	s.conn, s.jetstream = conn, js
	return nil
}

func (s *Sink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	for _, e := range envelopes {
		contents, err := e.AsByte()
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not marshal envelope")
			return err
		}
		_, err = s.jetstream.Publish(s.validSubject, contents)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not publish valid envelope to jetstream")
			return err
		}
	}
	return nil
}

func (s *Sink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	for _, e := range envelopes {
		contents, err := e.AsByte()
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not marshal envelope")
			return err
		}
		_, err = s.jetstream.Publish(s.invalidSubject, contents)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not publish invalid envelope to jetstream")
			return err
		}
	}
	return nil
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	return nil
}

func (s *Sink) Close() {
	log.Debug().Msg("🟡 closing nats sink")
	s.conn.Close()
}
