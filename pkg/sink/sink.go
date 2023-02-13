// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/db"
	"github.com/silverton-io/buz/pkg/envelope"
	"golang.org/x/net/context"
)

const (
	PUBSUB           string = "pubsub"
	REDPANDA         string = "redpanda"
	KAFKA            string = "kafka"
	KINESIS          string = "kinesis"
	KINESIS_FIREHOSE string = "kinesis-firehose"
	STDOUT           string = "stdout"
	HTTP             string = "http"
	HTTPS            string = "https"
	BLACKHOLE        string = "blackhole"
	FILE             string = "file"
	PUBNUB           string = "pubnub"
	NATS             string = "nats"
	NATS_JETSTREAM   string = "nats-jetstream"
	INDICATIVE       string = "indicative"
	AMPLITUDE        string = "amplitude"
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
	case PUBSUB:
		sink := PubsubSink{}
		return &sink, nil
	case KAFKA:
		sink := KafkaSink{}
		return &sink, nil
	case REDPANDA:
		sink := KafkaSink{}
		return &sink, nil
	case KINESIS:
		sink := KinesisSink{}
		return &sink, nil
	case KINESIS_FIREHOSE:
		sink := KinesisFirehoseSink{}
		return &sink, nil
	case STDOUT:
		sink := StdoutSink{}
		return &sink, nil
	case HTTP:
		sink := HttpSink{}
		return &sink, nil
	case HTTPS:
		sink := HttpSink{}
		return &sink, nil
	case db.ELASTICSEARCH:
		sink := ElasticsearchSink{}
		return &sink, nil
	case BLACKHOLE:
		sink := BlackholeSink{}
		return &sink, nil
	case FILE:
		sink := FileSink{}
		return &sink, nil
	case PUBNUB:
		sink := PubnubSink{}
		return &sink, nil
	case db.POSTGRES:
		sink := PostgresSink{}
		return &sink, nil
	case db.MYSQL:
		sink := MysqlSink{}
		return &sink, nil
	case db.MATERIALIZE:
		sink := MaterializeSink{}
		return &sink, nil
	case db.CLICKHOUSE:
		sink := ClickhouseSink{}
		return &sink, nil
	case db.MONGODB:
		sink := MongodbSink{}
		return &sink, nil
	case db.TIMESCALE:
		sink := TimescaleSink{}
		return &sink, nil
	case NATS:
		sink := NatsSink{}
		return &sink, nil
	case INDICATIVE:
		sink := IndicativeSink{}
		return &sink, nil
	case AMPLITUDE:
		sink := AmplitudeSink{}
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
