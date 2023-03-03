// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package nats

import (
	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	id                   *uuid.UUID
	sinkType             string
	name                 string
	deliveryRequired     bool
	conn                 *nats.Conn
	encodedConn          *nats.EncodedConn
	defaultEventsSubject string
	input                chan []envelope.Envelope
	shutdown             chan int
	// FIXME! Add .creds/token/tls cert/nkey auth
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		SinkType:         s.sinkType,
		DeliveryRequired: s.deliveryRequired,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing nats sink")
	id := uuid.New()
	s.id, s.sinkType, s.name, s.deliveryRequired = &id, conf.Type, conf.Name, conf.DeliveryRequired
	conn, err := nats.Connect(conf.Hosts[0], nats.UserInfo(conf.User, conf.Password))
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open nats connection")
		return err
	}
	encodedConn, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open encoded connection")
		return err
	}
	s.conn, s.encodedConn = conn, encodedConn
	s.defaultEventsSubject = constants.BUZ_EVENTS
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	return nil
}

func (s *Sink) StartWorker() error {
	err := backendutils.StartSinkWorker(s.input, s.shutdown, s)
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.input <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	for _, e := range envelopes {
		err := s.encodedConn.Publish(s.defaultEventsSubject, &e) // FIXME -> shard
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Sink) Shutdown() error {
	log.Debug().Msg("🟢 shutting down " + s.sinkType + " sink")
	s.shutdown <- 1
	s.conn.Close()
	s.encodedConn.Close()
	return nil
}
