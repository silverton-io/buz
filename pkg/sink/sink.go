package sink

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/tele"
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
)

type Sink interface {
	Initialize(conf config.Sink)
	BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope)
	BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope)
	BatchPublishValidAndInvalid(ctx context.Context, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope, meta *tele.Meta)
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
	default:
		e := errors.New("unsupported sink: " + conf.Type)
		log.Error().Stack().Err(e).Msg("unsupported sink")
		return nil, e
	}
}

func InitializeSink(conf config.Sink, s Sink) {
	s.Initialize(conf)
	log.Info().Msg(conf.Type + " sink initialized")
}
