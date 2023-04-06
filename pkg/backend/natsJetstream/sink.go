// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package natsJetstream

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink struct {
	metadata  backendutils.SinkMetadata
	conn      *nats.Conn
	jetstream nats.JetStreamContext
	// FIXME! Add .creds/token/tls cert/nkey auth
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing nats jetstream sink")
	s.metadata = backendutils.NewSinkMetadataFromConfig(conf)
	conn, err := nats.Connect(conf.Hosts[0], nats.UserInfo(conf.User, conf.Password))
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not open nats connection")
		return err
	}
	js, err := conn.JetStream()

	if err != nil {
		log.Error().Err(err).Msg("🔴 could not use jetstream context")
		return err
	}
	s.conn, s.jetstream = conn, js
	return nil
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	for _, e := range envelopes { // FIXME -> shard
		contents, err := e.AsByte()
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not marshal envelope")
			return err
		}
		_, err = s.jetstream.Publish(s.metadata.DefaultOutput, contents)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not publish valid envelope to jetstream")
			return err
		}
	}
	return nil
}

func (s *Sink) Close() {
	log.Debug().Msg("🟡 closing nats sink")
	s.conn.Close()
}
