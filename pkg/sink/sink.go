// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/backend/blackhole"
	"github.com/silverton-io/buz/pkg/backend/file"
	"github.com/silverton-io/buz/pkg/backend/stdout"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
)

type Sink interface {
	Metadata() backendutils.SinkMetadata
	Initialize(conf config.Sink) error
	Enqueue(envelopes []envelope.Envelope)
	Shutdown() error
}

func BuildSink(conf config.Sink) (sink Sink, err error) {
	switch conf.Type {
	// System
	case constants.BLACKHOLE:
		sink := blackhole.Sink{}
		return &sink, nil
	case constants.FILE:
		sink := file.Sink{}
		return &sink, nil
	case constants.STDOUT:
		sink := stdout.Sink{}
		return &sink, nil
	// Streams
	// case constants.PUBSUB:
	// 	sink := pubsub.Sink{}
	// 	return &sink, nil
	// case constants.KAFKA:
	// 	sink := kafka.Sink{}
	// 	return &sink, nil
	// case constants.REDPANDA:
	// 	sink := kafka.Sink{}
	// 	return &sink, nil
	// case constants.KINESIS:
	// 	sink := kinesis.Sink{}
	// 	return &sink, nil
	// case constants.KINESIS_FIREHOSE:
	// 	sink := kinesisFirehose.Sink{}
	// 	return &sink, nil
	// case constants.HTTP:
	// 	sink := http.Sink{}
	// 	return &sink, nil
	// case constants.HTTPS:
	// 	sink := http.Sink{}
	// 	return &sink, nil
	// case constants.ELASTICSEARCH:
	// 	sink := elasticsearch.Sink{}
	// 	return &sink, nil
	// case constants.PUBNUB:
	// 	sink := pubnub.Sink{}
	// 	return &sink, nil
	// case constants.POSTGRES:
	// 	sink := postgresdb.Sink{}
	// 	return &sink, nil
	// case constants.MYSQL:
	// 	sink := mysqldb.Sink{}
	// 	return &sink, nil
	// case constants.MATERIALIZE:
	// 	sink := materializedb.Sink{}
	// 	return &sink, nil
	// case constants.CLICKHOUSE:
	// 	sink := clickhousedb.Sink{}
	// 	return &sink, nil
	// case constants.MONGODB:
	// 	sink := mongodb.Sink{}
	// 	return &sink, nil
	// case constants.TIMESCALE:
	// 	sink := timescaledb.Sink{}
	// 	return &sink, nil
	// case constants.NATS:
	// 	sink := nats.Sink{}
	// 	return &sink, nil
	// case constants.AMPLITUDE:
	// 	sink := amplitude.Sink{}
	// 	return &sink, nil
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
