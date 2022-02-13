package forwarder

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"golang.org/x/net/context"
)

const (
	PUBSUB string = "pubsub"
	KAFKA  string = "kafka"
)

type Forwarder interface {
	Initialize(config config.Forwarder)
	PublishValid(ctx context.Context, event interface{})
	PublishInvalid(ctx context.Context, event interface{})
	BatchPublishValid(ctx context.Context, events []interface{})
	BatchPublishInvalid(ctx context.Context, events []interface{})
	Close()
}

func BuildForwarder(config config.Forwarder) (forwarder Forwarder, err error) {
	switch config.Type {
	case PUBSUB:
		forwarder := PubsubForwarder{}
		forwarder.Initialize(config)
		log.Debug().Msg("pubsub forwarder initialized")
		return &forwarder, nil
	case KAFKA:
		forwarder := KafkaForwarder{}
		forwarder.Initialize(config)
		log.Debug().Msg("kafka forwarder initialized")
		return &forwarder, nil
	default:
		e := errors.New("unsupported forwarder: " + config.Type)
		log.Fatal().Err(e).Msg("unsupported forwarder")
		return nil, e
	}
}
