package sink

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

type Sink interface {
	Initialize(config config.Forwarder)
	BatchPublishValid(ctx context.Context, events []interface{})
	BatchPublishInvalid(ctx context.Context, events []interface{})
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

func BatchPublishValidAndInvalid(inputType string, sink Sink, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta) {
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
	sink.BatchPublishValid(ctx, validEvents)
	sink.BatchPublishInvalid(ctx, invalidEvents)
	// Increment global metadata counters
	atomic.AddInt64(validCounter, int64(len(validEvents)))
	atomic.AddInt64(invalidCounter, int64(len(invalidEvents)))
}
