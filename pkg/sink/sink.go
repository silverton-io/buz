// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/backend/blackhole"
	"github.com/silverton-io/buz/pkg/backend/elasticsearch"
	"github.com/silverton-io/buz/pkg/backend/eventbridge"
	"github.com/silverton-io/buz/pkg/backend/file"
	"github.com/silverton-io/buz/pkg/backend/http"
	"github.com/silverton-io/buz/pkg/backend/kafka"
	"github.com/silverton-io/buz/pkg/backend/kinesis"
	"github.com/silverton-io/buz/pkg/backend/kinesisFirehose"
	"github.com/silverton-io/buz/pkg/backend/mongodb"
	"github.com/silverton-io/buz/pkg/backend/mysqldb"
	"github.com/silverton-io/buz/pkg/backend/nats"
	"github.com/silverton-io/buz/pkg/backend/postgresdb"
	"github.com/silverton-io/buz/pkg/backend/pubnub"
	"github.com/silverton-io/buz/pkg/backend/pubsub"
	"github.com/silverton-io/buz/pkg/backend/stdout"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
)

func getSink(conf config.Sink) (sink backendutils.Sink, err error) {
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
	case constants.HTTP:
		sink := http.Sink{}
		return &sink, nil
	case constants.HTTPS:
		sink := http.Sink{}
		return &sink, nil
	// Streams
	case constants.KAFKA:
		sink := kafka.Sink{}
		return &sink, nil
	case constants.REDPANDA:
		sink := kafka.Sink{}
		return &sink, nil
	case constants.PUBSUB:
		sink := pubsub.Sink{}
		return &sink, nil
	case constants.KINESIS:
		sink := kinesis.Sink{}
		return &sink, nil
	case constants.KINESIS_FIREHOSE:
		sink := kinesisFirehose.Sink{}
		return &sink, nil
	case constants.EVENTBRIDGE:
		sink := eventbridge.Sink{}
		return &sink, nil
	// Message Brokers
	case constants.NATS:
		sink := nats.Sink{}
		return &sink, nil
	// Databases
	case constants.POSTGRES:
		sink := postgresdb.Sink{}
		return &sink, nil
	case constants.TIMESCALE:
		sink := postgresdb.Sink{}
		return &sink, nil
	case constants.MYSQL:
		sink := mysqldb.Sink{}
		return &sink, nil
	case constants.MONGODB:
		sink := mongodb.Sink{}
		return &sink, nil
	case constants.ELASTICSEARCH:
		sink := elasticsearch.Sink{}
		return &sink, nil
	// case constants.CLICKHOUSE:
	// 	sink := clickhousedb.Sink{}
	// 	return &sink, nil
	case constants.PUBNUB:
		sink := pubnub.Sink{}
		return &sink, nil
	case constants.MATERIALIZE:
		sink := postgresdb.Sink{}
		return &sink, nil
	default:
		e := errors.New("unsupported sink: " + conf.Type)
		log.Error().Stack().Err(e).Msg("🔴 unsupported sink")
		return nil, e
	}
}

func NewSink(conf config.Sink) (backendutils.Sink, error) {
	sink, _ := getSink(conf)
	err := sink.Initialize(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("🔴 could not initialize sink")
		return nil, err
	}
	err = sink.StartWorker()
	if err != nil {
		log.Fatal().Err(err).Interface("metadata", sink.Metadata()).Msg("could not start sink worker")
	}
	log.Info().Msg("🟢 " + conf.Type + " sink initialized")
	return sink, nil
}

func BuildAndInitializeSinks(conf []config.Sink) ([]backendutils.Sink, error) {
	var sinks []backendutils.Sink
	for _, sConf := range conf {
		sink, err := NewSink(sConf)
		if err != nil {
			log.Fatal().Err(err).Msg("🔴 could not initialize sink")
			return nil, err
		}
		sinks = append(sinks, sink)
	}
	return sinks, nil
}
