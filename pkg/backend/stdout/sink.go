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
	sinkType         string
	name             string
	deliveryRequired bool
	input            chan []envelope.Envelope
	shutdown         chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		SinkType:         s.sinkType,
		DeliveryRequired: s.deliveryRequired,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing stdout sink")
	id := uuid.New()
	s.id, s.name, s.sinkType = &id, conf.Name, conf.Type
	s.deliveryRequired = conf.DeliveryRequired
	s.input = make(chan []envelope.Envelope, 10000)
	s.shutdown = make(chan int, 1)
	return nil
}

func (s *Sink) StartWorker() error {
	err := backendutils.StartSinkWorker(s.input, s.shutdown, s)
	return err
}

func (s *Sink) Enqueue(envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.input <- envelopes
	return nil
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	var validEnvelopes []envelope.Envelope
	var invalidEnvelopes []envelope.Envelope
	for _, e := range envelopes {
		if e.IsValid {
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
	s.shutdown <- 1
	return nil
}
