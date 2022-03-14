package sink

import (
	"errors"
	"sync/atomic"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/input"
	"github.com/silverton-io/honeypot/pkg/tele"
	"golang.org/x/net/context"
)

const (
	PUBSUB           string = "pubsub"
	KAFKA            string = "kafka"
	KINESIS          string = "kinesis"
	KINESIS_FIREHOSE string = "kinesis-firehose"
	STDOUT           string = "stdout"
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

func incrementStats(inputType string, validCount int, invalidCount int, meta *tele.Meta) {
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
	atomic.AddInt64(validCounter, int64(validCount))
	atomic.AddInt64(invalidCounter, int64(invalidCount))
}
