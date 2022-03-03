package sink

import (
	"context"
	"encoding/json"
	"sync"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/tele"
)

type KinesisFirehoseSink struct {
	client              *firehose.Client
	validEventsStream   string
	invalidEventsStream string
}

func (s *KinesisFirehoseSink) Initialize(conf config.Sink) {
	ctx := context.Background()
	cfg, _ := awsconf.LoadDefaultConfig(ctx)
	client := firehose.NewFromConfig(cfg)
	s.client, s.validEventsStream, s.invalidEventsStream = client, conf.ValidEventTopic, conf.InvalidEventTopic
}

func (s *KinesisFirehoseSink) batchPublish(ctx context.Context, stream string, events []interface{}) {
	var wg sync.WaitGroup
	for _, event := range events {
		payload, _ := json.Marshal(event)
		newline := []byte("\n")
		payload = append(payload, newline...)
		record := types.Record{
			Data: payload,
		}
		input := &firehose.PutRecordInput{
			DeliveryStreamName: &stream,
			Record:             &record,
		}
		go func() {
			wg.Add(1)
			output, err := s.client.PutRecord(ctx, input) // Will want to use `PutRecordBatch`
			defer wg.Done()
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not publish event to kinesis firehose")
			} else {
				log.Debug().Msgf("published event " + *output.RecordId + " to stream " + stream)
			}
		}()
	}
	wg.Wait()
}

func (s *KinesisFirehoseSink) batchPublishValid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.validEventsStream, events)
}

func (s *KinesisFirehoseSink) batchPublishInvalid(ctx context.Context, events []interface{}) {
	s.batchPublish(ctx, s.invalidEventsStream, events)
}

func (s *KinesisFirehoseSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []interface{}, invalidEvents []interface{}, meta *tele.Meta) {
	// Publish
	s.batchPublishValid(ctx, validEvents)
	s.batchPublishInvalid(ctx, invalidEvents)
	// Increment stats counters
	incrementStats(inputType, len(validEvents), len(invalidEvents), meta)
}

func (s *KinesisFirehoseSink) Close() {
	log.Debug().Msg("closing kinesis firehose sink client")
}
