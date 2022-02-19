package forwarder

import (
	"errors"
	"sync/atomic"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"golang.org/x/net/context"
)

const (
	PUBSUB string = "pubsub"
	KAFKA  string = "kafka"
)

type Forwarder interface {
	Initialize(config config.Forwarder)
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

func BatchPublishValidAndInvalid(inputType string, forwarder Forwarder, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta) {
	ctx := context.Background()
	var validCounter *int64
	var invalidCounter *int64
	switch inputType {
	case input.GENERIC_INPUT:
		validCounter = &meta.ValidGenericEventsProcessed
		invalidCounter = &meta.InvalidGenericEventsProcessed
	case input.CLOUDEVENTS_INPUT:
		validCounter = &meta.ValidCloudEventsProcessed
		invalidCounter = &meta.InvalidCloudEventsProcessed
	default:
		validCounter = &meta.ValidSnowplowEventsProcessed
		invalidCounter = &meta.InvalidSnowplowEventsProcessed
	}
	// Publish
	forwarder.BatchPublishValid(ctx, validEvents)
	forwarder.BatchPublishInvalid(ctx, invalidEvents)
	// Increment global metadata counters
	atomic.AddInt64(validCounter, int64(len(validEvents)))
	atomic.AddInt64(invalidCounter, int64(len(invalidEvents)))
}
