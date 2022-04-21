package sink

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"golang.org/x/net/context"
)

const (
	PUBSUB           string = "pubsub"
	KAFKA            string = "kafka"
	KINESIS          string = "kinesis"
	KINESIS_FIREHOSE string = "kinesis-firehose"
	STDOUT           string = "stdout"
	HTTP             string = "http"
	HTTPS            string = "https"
	RELAY            string = "relay"
	ELASTICSEARCH    string = "elasticsearch"
	BLACKHOLE        string = "blackhole"
	FILE             string = "file"
	POSTGRES         string = "postgres"
	MYSQL            string = "mysql"
	MATERIALIZE      string = "materialize"
	CLICKHOUSE       string = "clickhouse"
	PUBNUB           string = "pubnub"
	MONGODB          string = "mongodb"
)

type Sink interface {
	Id() *uuid.UUID
	Name() string
	Initialize(conf config.Sink) error
	BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error
	BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error
	Close()
}

func BuildSink(conf config.Sink) (sink Sink, err error) {
	switch conf.Type {
	case PUBSUB:
		sink := PubsubSink{}
		return &sink, nil
	case KAFKA:
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
	case RELAY:
		sink := RelaySink{}
		return &sink, nil
	case ELASTICSEARCH:
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
	case POSTGRES:
		sink := PostgresSink{}
		return &sink, nil
	case MYSQL:
		sink := MysqlSink{}
		return &sink, nil
	case MATERIALIZE:
		sink := MaterializeSink{}
		return &sink, nil
	case CLICKHOUSE:
		sink := ClickhouseSink{}
		return &sink, nil
	case MONGODB:
		sink := MongodbSink{}
		return &sink, nil
	default:
		e := errors.New("unsupported sink: " + conf.Type)
		log.Error().Stack().Err(e).Msg("unsupported sink")
		return nil, e
	}
}

func InitializeSink(conf config.Sink, s Sink) error {
	err := s.Initialize(conf)
	if err != nil {
		log.Debug().Stack().Err(err).Msg("could not initialize sink")
		return err
	}
	log.Info().Msg(conf.Type + " sink initialized")
	return nil
}

func BuildAndInitializeSinks(conf []config.Sink) ([]Sink, error) {
	var sinks []Sink
	for _, sConf := range conf {
		sink, err := BuildSink(sConf)
		if err != nil {
			log.Debug().Err(err).Msg("could not build sink")
			return nil, err
		}
		err = InitializeSink(sConf, sink)
		if err != nil {
			log.Debug().Err(err).Msg("could not initialize sink")
			return nil, err
		}
		sinks = append(sinks, sink)
	}
	return sinks, nil
}
