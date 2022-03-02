package sink

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/tele"
)

type KinesisSink struct {
	client              *kinesis.Client
	validEventsStream   string
	invalidEventsStream string
}

func (s *KinesisSink) Initialize(conf config.Sink) {
	ctx := context.Background()
	cfg, _ := awsconf.LoadDefaultConfig(ctx)
	client := kinesis.NewFromConfig(cfg)
	s.client, s.validEventsStream, s.invalidEventsStream = client, conf.ValidEventTopic, conf.InvalidEventTopic
}

func (s *KinesisSink) batchPublish(ctx context.Context, stream string, events []interface{}) {
	var wg sync.WaitGroup
	go func() {
		for _, event := range events {
			payload, _ := json.Marshal(event)
			input := &kinesis.PutRecordInput{
				Data: payload,
				// PartitionKey: partitionName, FIXME! Add configurable partition key
				StreamName: &stream,
			}
			wg.Add(1)
			_, err := s.client.PutRecord(ctx, input)
			defer wg.Done()
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not publish event to kinesis")
			} else {
				log.Debug().Msgf("published event to stream " + stream)
			}
		}
	}()
	wg.Wait()
}

func (s *KinesisSink) batchPublishValid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.validEventsStream, events)
}

func (s *KinesisSink) batchPublishInvalid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.invalidEventsStream, events)
}

func (s *KinesisSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta) {
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
	s.batchPublishValid(ctx, validEvents)
	s.batchPublishInvalid(ctx, invalidEvents)
	// Increment global metadata counters
	atomic.AddInt64(validCounter, int64(len(validEvents)))
	atomic.AddInt64(invalidCounter, int64(len(invalidEvents)))
}

func (s *KinesisSink) Close() {
	log.Debug().Msg("closing kinesis sink client")

}
