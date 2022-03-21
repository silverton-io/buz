package sink

import (
	"context"
	"encoding/json"
	"sync"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
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

func (s *KinesisSink) batchPublish(ctx context.Context, stream string, envelopes []envelope.Envelope) {
	var wg sync.WaitGroup
	for _, event := range envelopes {
		partitionKey := "blah" // FIXME!
		payload, _ := json.Marshal(event)
		input := &kinesis.PutRecordInput{
			Data:         payload,
			PartitionKey: &partitionKey,
			StreamName:   &stream,
		}
		wg.Add(1)
		go func() {
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

func (s *KinesisSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.validEventsStream, envelopes)
}

func (s *KinesisSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.invalidEventsStream, envelopes)
}

func (s *KinesisSink) BatchPublishValidAndInvalid(ctx context.Context, validEnvelopes []envelope.Envelope, invalidEnvelopes []envelope.Envelope, meta *tele.Meta) {
	// Publish
	go s.BatchPublishValid(ctx, validEnvelopes)
	go s.BatchPublishInvalid(ctx, invalidEnvelopes)
	// Increment stats counters
}

func (s *KinesisSink) Close() {
	log.Debug().Msg("closing kinesis sink client")

}
