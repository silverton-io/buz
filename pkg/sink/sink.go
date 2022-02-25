package sink

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"golang.org/x/net/context"
)

const (
	PUBSUB string = "pubsub"
	KAFKA  string = "kafka"
)

type Sink interface {
	Initialize(config config.Forwarder)
	BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta)
	Close()
}

func BuildSink(config config.Forwarder) (sink Sink, err error) {
	switch config.Type {
	case PUBSUB:
		sink := PubsubSink{}
		sink.Initialize(config)
		log.Debug().Msg("pubsub sink initialized")
		return &sink, nil
	case KAFKA:
		sink := KafkaSink{}
		sink.Initialize(config)
		log.Debug().Msg("kafka sink initialized")
		return &sink, nil
	default:
		e := errors.New("unsupported sink: " + config.Type)
		log.Fatal().Err(e).Msg("unsupported sink")
		return nil, e
	}
}
