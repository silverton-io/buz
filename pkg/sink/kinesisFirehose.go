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
	"github.com/silverton-io/honeypot/pkg/envelope"
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

func (s *KinesisFirehoseSink) batchPublish(ctx context.Context, stream string, envelopes []envelope.Envelope) {
	var wg sync.WaitGroup
	for _, event := range envelopes {
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
		wg.Add(1)
		go func() {
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

func (s *KinesisFirehoseSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.validEventsStream, envelopes)
}

func (s *KinesisFirehoseSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.invalidEventsStream, envelopes)
}

func (s *KinesisFirehoseSink) BatchPublishValidAndInvalid(ctx context.Context, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope) {
	// Publish
	go s.BatchPublishValid(ctx, validEnvelopes)
	go s.BatchPublishInvalid(ctx, invalidEnvelopes)
	// Increment stats counters
}

func (s *KinesisFirehoseSink) Close() {
	log.Debug().Msg("closing kinesis firehose sink client")
}
