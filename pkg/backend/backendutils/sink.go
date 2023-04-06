// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package backendutils

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

var DEFAULT_SINK_TIMEOUT_SECONDS int = 15

type SinkMetadata struct {
	Id               uuid.UUID `json:"id"`
	SinkType         string    `json:"sinkType"`
	Name             string    `json:"name"`
	DeliveryRequired bool      `json:"deliveryRequired"`
	DefaultOutput    string    `json:"defaultOutput"`
	DeadletterOutput string    `json:"deadletterOutput"`
}

func NewSinkMetadataFromConfig(conf config.Sink) SinkMetadata {
	return SinkMetadata{
		Id:               uuid.New(),
		SinkType:         conf.Type,
		Name:             conf.Name,
		DeliveryRequired: conf.DeliveryRequired,
		DefaultOutput:    conf.DefaultOutput,
		DeadletterOutput: conf.DeadletterOutput,
	}
}

type Sink interface {
	Metadata() SinkMetadata
	Initialize(conf config.Sink) error
	StartWorker() error
	Enqueue(envelopes []envelope.Envelope) error
	Dequeue(ctx context.Context, envelopes []envelope.Envelope, output string) error
	Shutdown() error
}

func publish(ctx context.Context, sink Sink, envelopes []envelope.Envelope, output string) error {
	if len(envelopes) > 0 {
		err := sink.Dequeue(ctx, envelopes, output)
		if err != nil {
			log.Error().Err(err).Interface("metadata", sink.Metadata()).Msg("could not dequeue envelopes to output " + output)
		}
	}
	return nil
}

// Each sink runs an associated worker goroutine, which is responsible
// for dequeuing envelopes.
func StartSinkWorker(input <-chan []envelope.Envelope, shutdown <-chan int, sink Sink) error {
	go func(input <-chan []envelope.Envelope, shutdown <-chan int, sink Sink) {
		for {
			select {
			case envelopes := <-input:
				// Just handle valid/invalid for now. This will be where events will be further sharded going forward.
				var invalidEnvelopes []envelope.Envelope
				var validEnvelopes []envelope.Envelope
				for _, envelope := range envelopes {
					if envelope.IsValid {
						validEnvelopes = append(validEnvelopes, envelope)
					} else {
						invalidEnvelopes = append(invalidEnvelopes, envelope)
					}
				}
				ctx := context.Background()
				// Send good events along
				err := publish(ctx, sink, validEnvelopes, sink.Metadata().DefaultOutput)
				if err != nil {
					log.Error().Err(err).Interface("metadata", sink.Metadata()).Msg("could not publish envelopes to sink")
				}
				// Send bad events to deadletter
				err = publish(ctx, sink, invalidEnvelopes, sink.Metadata().DeadletterOutput)
				if err != nil {
					log.Error().Err(err).Interface("metadata", sink.Metadata()).Msg("could not publish envelopes to sink")
				}
			case <-shutdown:
				return
			}
		}
	}(input, shutdown, sink)
	return nil
}
