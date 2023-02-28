// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package backendutils

import (
	"context"

	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
)

type SinkMetadata struct {
	Id               *uuid.UUID `json:"id"`
	SinkType         string     `json:"sinkType"`
	Name             string     `json:"name"`
	DeliveryRequired bool       `json:"deliveryRequired"`
}
type Sink interface {
	Metadata() SinkMetadata
	Initialize(conf config.Sink) error
	StartWorker() error
	Enqueue(envelopes []envelope.Envelope) error
	Dequeue(ctx context.Context, envelopes []envelope.Envelope) error
	Shutdown() error
}

// Each sink runs an associated worker goroutine, which is responsible
// for dequeuing envelopes.
func StartSinkWorker(input <-chan []envelope.Envelope, shutdown <-chan int, sink Sink) error {
	go func(input <-chan []envelope.Envelope, shutdown <-chan int, sink Sink) {
		for {
			select {
			case envelopes := <-input:
				ctx := context.Background()
				sink.Dequeue(ctx, envelopes)
			case <-shutdown:
				return
			}
		}
	}(input, shutdown, sink)
	return nil
}
