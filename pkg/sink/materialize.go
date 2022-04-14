package sink

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type MaterializeSink struct {
	id           *uuid.UUID
	name         string
	conn         *pgx.Conn
	validTable   string
	invalidTable string
}

func (s *MaterializeSink) Id() *uuid.UUID {
	return s.id
}

func (s *MaterializeSink) Name() string {
	return s.name
}

func (s *MaterializeSink) Initialize(conf config.Sink) {
	log.Debug().Msg("intializing materialize sink")
	s.validTable, s.invalidTable = conf.ValidTable, conf.InvalidTable
}

func (s *MaterializeSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {

}

func (s *MaterializeSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {

}

func (s *MaterializeSink) Close() {
	log.Debug().Msg("closing materialize sink")
	s.conn.Close()
}
