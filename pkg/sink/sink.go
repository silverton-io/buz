// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/amplitude"
	"github.com/silverton-io/buz/pkg/backend/blackhole"
	"github.com/silverton-io/buz/pkg/backend/clickhousedb"
	"github.com/silverton-io/buz/pkg/backend/elasticsearch"
	"github.com/silverton-io/buz/pkg/backend/file"
	"github.com/silverton-io/buz/pkg/backend/http"
	"github.com/silverton-io/buz/pkg/backend/kafka"
	"github.com/silverton-io/buz/pkg/backend/materializedb"
	"github.com/silverton-io/buz/pkg/backend/mongodb"
	"github.com/silverton-io/buz/pkg/backend/mysqldb"
	"github.com/silverton-io/buz/pkg/backend/postgresdb"
	"github.com/silverton-io/buz/pkg/backend/stdout"
	"github.com/silverton-io/buz/pkg/backend/timescaledb"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
	"golang.org/x/net/context"
)

type Sink interface {
	Id() *uuid.UUID
	Name() string
	Type() string
	DeliveryRequired() bool
	Initialize(conf config.Sink) error
	BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error
	Close()
	// FIXME! Add "shard by" mechanism.
}

func BuildSink(conf config.Sink) (sink Sink, err error) {
	switch conf.Type {
	case constants.PUBSUB:
		sink := PubsubSink{}
		return &sink, nil
	case constants.KAFKA:
		sink := kafka.Sink{}
		return &sink, nil
	case constants.REDPANDA:
		sink := kafka.Sink{}
		return &sink, nil
	case constants.KINESIS:
		sink := KinesisSink{}
		return &sink, nil
	case constants.KINESIS_FIREHOSE:
		sink := KinesisFirehoseSink{}
		return &sink, nil
	case constants.STDOUT:
		sink := stdout.Sink{}
		return &sink, nil
	case constants.HTTP:
		sink := http.Sink{}
		return &sink, nil
	case constants.HTTPS:
		sink := http.Sink{}
		return &sink, nil
	case constants.ELASTICSEARCH:
		sink := elasticsearch.Sink{}
		return &sink, nil
	case constants.BLACKHOLE:
		sink := blackhole.Sink{}
		return &sink, nil
	case constants.FILE:
		sink := file.Sink{}
		return &sink, nil
	case constants.PUBNUB:
		sink := PubnubSink{}
		return &sink, nil
	case constants.POSTGRES:
		sink := postgresdb.Sink{}
		return &sink, nil
	case constants.MYSQL:
		sink := mysqldb.Sink{}
		return &sink, nil
	case constants.MATERIALIZE:
		sink := materializedb.Sink{}
		return &sink, nil
	case constants.CLICKHOUSE:
		sink := clickhousedb.Sink{}
		return &sink, nil
	case constants.MONGODB:
		sink := mongodb.Sink{}
		return &sink, nil
	case constants.TIMESCALE:
		sink := timescaledb.Sink{}
		return &sink, nil
	case constants.NATS:
		sink := NatsSink{}
		return &sink, nil
	case constants.AMPLITUDE:
		sink := amplitude.Sink{}
		return &sink, nil
	default:
		e := errors.New("unsupported sink: " + conf.Type)
		log.Error().Stack().Err(e).Msg("🔴 unsupported sink")
		return nil, e
	}
}

func InitializeSink(conf config.Sink, s Sink) error {
	err := s.Initialize(conf)
	if err != nil {
		log.Error().Err(err).Msg("🔴 could not initialize sink")
		return err
	}
	log.Info().Msg("🟢 " + conf.Type + " sink initialized")
	return nil
}

func BuildAndInitializeSinks(conf []config.Sink) ([]Sink, error) {
	var sinks []Sink
	for _, sConf := range conf {
		sink, err := BuildSink(sConf)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not build sink")
			return nil, err
		}
		err = InitializeSink(sConf, sink)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not initialize sink")
			return nil, err
		}
		sinks = append(sinks, sink)
	}
	return sinks, nil
}
