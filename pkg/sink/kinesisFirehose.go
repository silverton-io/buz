package sink

import (
	"context"
	"encoding/json"
	"sync"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

type KinesisFirehoseSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	client           *firehose.Client
	validStream      string
	invalidStream    string
}

func (s *KinesisFirehoseSink) Id() *uuid.UUID {
	return s.id
}

func (s *KinesisFirehoseSink) Name() string {
	return s.name
}

func (s *KinesisFirehoseSink) Type() string {
	return KINESIS_FIREHOSE
}

func (s *KinesisFirehoseSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *KinesisFirehoseSink) Initialize(conf config.Sink) error {
	ctx := context.Background()
	cfg, err := awsconf.LoadDefaultConfig(ctx)
	client := firehose.NewFromConfig(cfg)
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.client, s.validStream, s.invalidStream = client, conf.ValidStream, conf.InvalidStream
	return err
}

func (s *KinesisFirehoseSink) batchPublish(ctx context.Context, stream string, envelopes []envelope.Envelope) error {
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
		pubErr := make(chan error, 1)
		go func(pErr chan error) {
			output, err := s.client.PutRecord(ctx, input) // Will want to use `PutRecordBatch`
			defer wg.Done()
			if err != nil {
				log.Error().Stack().Err(err).Msg("could not publish event to kinesis firehose")
				pubErr <- err
			} else {
				log.Debug().Msgf("published event " + *output.RecordId + " to stream " + stream)
				pubErr <- nil
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

func (s *KinesisFirehoseSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.validStream, envelopes)
	return err
}

func (s *KinesisFirehoseSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, s.invalidStream, envelopes)
	return err
}

func (s *KinesisFirehoseSink) Close() {
	log.Debug().Msg("closing kinesis firehose sink client")
}
