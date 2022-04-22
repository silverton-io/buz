package sink

import (
	"context"
	"encoding/json"
	"sync"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type KinesisSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	client           *kinesis.Client
	validStream      string
	invalidStream    string
}

func (s *KinesisSink) Id() *uuid.UUID {
	return s.id
}

func (s *KinesisSink) Name() string {
	return s.name
}

func (s *KinesisSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *KinesisSink) Initialize(conf config.Sink) error {
	ctx := context.Background()
	cfg, err := awsconf.LoadDefaultConfig(ctx)
	client := kinesis.NewFromConfig(cfg)
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.client, s.validStream, s.invalidStream = client, conf.ValidStream, conf.InvalidStream
	return err
}

func (s *KinesisSink) batchPublish(ctx context.Context, stream string, envelopes []envelope.Envelope) error {
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
		pubErr := make(chan error, 1)
		go func(pErr chan error) {
			output, err := s.client.PutRecord(ctx, input) // Will want to use `PutRecordBatch`
			defer wg.Done()
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not publish event to kinesis")
				pErr <- err
			} else {
				log.Debug().Msgf("published event " + *output.SequenceNumber + " to stream " + stream)
				pErr <- nil
			}
		}(pubErr)
		err := <-pubErr
		if err != nil {
			return err
		}
	}
	wg.Wait()
	return nil
}

func (s *KinesisSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validStream, envelopes)
	return err
}

func (s *KinesisSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.invalidStream, envelopes)
	return err
}

func (s *KinesisSink) Close() {
	log.Debug().Msg("closing kinesis sink client")

}
