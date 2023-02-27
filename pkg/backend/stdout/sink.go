// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package stdout

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/util"
)

var (
	Info = Teal
	Warn = Yellow
	Fata = Red
)

var (
	Black   = Colorize("\033[1;30m%s\033[0m")
	Red     = Colorize("\033[1;31m%s\033[0m")
	Green   = Colorize("\033[1;32m%s\033[0m")
	Yellow  = Colorize("\033[1;33m%s\033[0m")
	Purple  = Colorize("\033[1;34m%s\033[0m")
	Magenta = Colorize("\033[1;35m%s\033[0m")
	Teal    = Colorize("\033[1;36m%s\033[0m")
	White   = Colorize("\033[1;37m%s\033[0m")
)

func Colorize(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	inputChan        chan []envelope.Envelope
	shutdownChan     chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	sinkType := "stdout"
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		Type:             sinkType,
		DeliveryRequired: s.deliveryRequired,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing stdout sink")
	id := uuid.New()
	s.id, s.name = &id, conf.Name
	s.deliveryRequired = conf.DeliveryRequired
	s.inputChan = make(chan []envelope.Envelope, 10000)
	s.shutdownChan = make(chan int, 1)
	go func(s *Sink) {
		for {
			select {
			case envelopes := <-s.inputChan:
				ctx := context.Background()
				s.Dequeue(ctx, envelopes)
			case <-s.shutdownChan:
				err := s.Shutdown()
				if err != nil {
					log.Error().Err(err).Interface("metadata", s.Metadata()).Msg("sink did not safely shut down")
				}
				return
			}
		}
	}(s)
	return nil
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.inputChan <- envelopes
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	var validEnvelopes []envelope.Envelope
	var invalidEnvelopes []envelope.Envelope
	for _, e := range envelopes {
		if *e.IsValid {
			validEnvelopes = append(validEnvelopes, e)
		} else {
			invalidEnvelopes = append(invalidEnvelopes, e)
		}
	}
	if len(validEnvelopes) > 0 {
		validEnvelopes := util.Stringify(validEnvelopes)
		fmt.Println(Green(validEnvelopes))
	}
	if len(invalidEnvelopes) > 0 {
		invalidEnvelopes := util.Stringify(invalidEnvelopes)
		fmt.Println(Red(invalidEnvelopes))
	}
	return nil
}

func (s *Sink) Shutdown() error {
	log.Debug().Msg("🟢 shutting down stdout sink")
	return nil
}
