// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type NatsJetstreamSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	conn             *nats.Conn
	jetstream        nats.JetStreamContext
	validSubject     string
	invalidSubject   string
	// FIXME! Add .creds/token/tls cert/nkey auth
}

func (s *NatsJetstreamSink) Id() *uuid.UUID {
	return s.id
}

func (s *NatsJetstreamSink) Name() string {
	return s.name
}

func (s *NatsJetstreamSink) Type() string {
	return NATS_JETSTREAM
}

func (s *NatsJetstreamSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *NatsJetstreamSink) Initialize(conf config.Sink) error {
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

	s.validSubject, s.invalidSubject = conf.ValidSubject, conf.InvalidSubject
	s.conn, s.jetstream = conn, js
	return nil
}

func (s *NatsJetstreamSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
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

func (s *NatsJetstreamSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
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

func (s *NatsJetstreamSink) Close() {
	log.Debug().Msg("🟡 closing nats sink")
	s.conn.Close()
}
