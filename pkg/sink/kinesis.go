package sink

import (
	"context"
	"encoding/json"
	"sync"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/tele"
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

func (s *KinesisSink) batchPublish(ctx context.Context, stream string, envelopes []event.Envelope) {
	var wg sync.WaitGroup
	for _, event := range envelopes {
		partitionKey := "blah"
		payload, _ := json.Marshal(event)
		input := &kinesis.PutRecordInput{
			Data:         payload,
			PartitionKey: &partitionKey,
			StreamName:   &stream,
		}
		go func() {
			wg.Add(1)
			output, err := s.client.PutRecord(ctx, input) // Will want to use `PutRecordBatch`
			defer wg.Done()
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not publish event to kinesis")
			} else {
				log.Debug().Msgf("published event " + *output.SequenceNumber + " to stream " + stream)
			}
		}()
	}
	wg.Wait()
}

func (s *KinesisSink) batchPublishValid(ctx context.Context, envelopes []event.Envelope) {
	s.batchPublish(ctx, s.validEventsStream, envelopes)
}

func (s *KinesisSink) batchPublishInvalid(ctx context.Context, envelopes []event.Envelope) {
	s.batchPublish(ctx, s.invalidEventsStream, envelopes)
}

func (s *KinesisSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEnvelopes []event.Envelope, invalidEnvelopes []event.Envelope, meta *tele.Meta) {
	// Publish
	go s.batchPublishValid(ctx, validEnvelopes)
	go s.batchPublishInvalid(ctx, invalidEnvelopes)
	// Increment stats counters
	incrementStats(inputType, len(validEnvelopes), len(invalidEnvelopes), meta)
}

func (s *KinesisSink) Close() {
	log.Debug().Msg("closing kinesis sink client")

}
