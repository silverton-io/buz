package sink

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"golang.org/x/net/context"
)

const (
	PUBSUB  string = "pubsub"
	KAFKA   string = "kafka"
	KINESIS string = "kinesis"
)

type Sink interface {
	Initialize(conf config.Sink)
	BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta)
	Close()
}

func BuildSink(conf config.Sink) (sink Sink, err error) {
	switch conf.Type {
	case PUBSUB:
		sink := PubsubSink{}
		sink.Initialize(conf)
		log.Debug().Msg("pubsub sink initialized")
		return &sink, nil
	case KAFKA:
		sink := KafkaSink{}
		sink.Initialize(conf)
		log.Debug().Msg("kafka sink initialized")
		return &sink, nil
	case KINESIS:
		sink := KinesisSink{}
		sink.Initialize(conf)
		log.Debug().Msg("kinesis sink initialized")
		return &sink, nil
	default:
		e := errors.New("unsupported sink: " + conf.Type)
		log.Fatal().Err(e).Msg("unsupported sink")
		return nil, e
	}
}
