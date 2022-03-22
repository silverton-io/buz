package sink

import (
	"context"
	"database/sql"

	"entgo.io/ent"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type ValidEnvelope struct {
	ent.Schema
}

func (ValidEnvelope) Fields() []ent.Field {
	return nil
}

type InvalidEnvelope struct {
	ent.Schema
}

func (InvalidEnvelope) Fields() []ent.Field {
	return nil
}

type PostgresSink struct {
	db *sql.DB
}

func (s *PostgresSink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing postgres sink")
	dbUrl := "postgresql://" + conf.DbUser + ":" + conf.DbPass + "@" + conf.DbHost + ":" + conf.DbPort + "/" + conf.DbName
	db, err := sql.Open("pgx", dbUrl)
	// drv := entsql.OpenDB(dialect.Postgres, db)
	// return ent.NewClient(ent.Driver(drv))
}

func (s *PostgresSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {}

func (s *PostgresSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
}

func (s *PostgresSink) Close() {
	log.Debug().Msg("closing postgres sink")
}
