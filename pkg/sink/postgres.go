package sink

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

func generateCreateTableSql(tableName string) string {
	return "create table if not exists " + tableName + "(id uuid, \"eventProtocol\" text, \"eventSchema\" text, source text, tstamp timestamp with time zone, ip text, \"isValid\" boolean, \"isRelayed\" boolean, \"validationError\" jsonb, payload jsonb);"
}

type PostgresSink struct {
	id           *uuid.UUID
	name         string
	conn         *pgx.Conn
	validTable   string
	invalidTable string
}

func (s *PostgresSink) Id() *uuid.UUID {
	return s.id
}

func (s *PostgresSink) Name() string {
	return s.name
}

func (s *PostgresSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing postgres sink")
	connectionConf := pgx.ConnConfig{
		Host:     conf.DbHost,
		Port:     conf.DbPort,
		Database: conf.DbName,
		User:     conf.DbUser,
		Password: conf.DbPass,
	}
	conn, err := pgx.Connect(connectionConf)
	if err != nil {
		log.Debug().Err(err).Msg("could not open db connection")
		return err
	}
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	s.conn, s.validTable, s.invalidTable = conn, conf.ValidTable, conf.InvalidTable
	createValidSql, createInvalidSql := generateCreateTableSql(s.validTable), generateCreateTableSql(s.invalidTable)
	for _, sql := range []string{createValidSql, createInvalidSql} {
		_, err := s.conn.Exec(sql)
		if err != nil {
			log.Debug().Err(err).Msg("could not create table")
			return err
		}
	}
	return nil
}

func (s *PostgresSink) batchPublish(ctx context.Context, tableName string, envelopes []envelope.Envelope) {
	var rows [][]interface{}
	for _, envelope := range envelopes {
		row := []interface{}{
			envelope.Id,
			envelope.EventProtocol,
			envelope.EventSchema,
			envelope.Source,
			envelope.Tstamp,
			envelope.Ip,
			envelope.IsValid,
			envelope.IsRelayed,
			envelope.ValidationError,
			envelope.Payload,
		}
		rows = append(rows, row)
	}
	copyCount, err := s.conn.CopyFrom(
		pgx.Identifier{tableName},
		[]string{"id", "eventProtocol", "eventSchema", "source", "tstamp", "ip", "isValid", "isRelayed", "validationError", "payload"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Error().Err(err).Msg("could not copy rows")
	}
	count := strconv.Itoa(copyCount)
	log.Debug().Msg("copied " + count + " rows to " + tableName)
}

func (s *PostgresSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.validTable, envelopes)
}

func (s *PostgresSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.invalidTable, envelopes)
}

func (s *PostgresSink) Close() {
	log.Debug().Msg("closing postgres sink")
	s.conn.Close()
}
