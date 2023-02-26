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
	fanout           bool
	inputChan        chan []envelope.Envelope
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	sinkType := "stdout"
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		Type:             sinkType,
		DeliveryRequired: s.deliveryRequired,
		Fanout:           false,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing stdout sink")
	id := uuid.New()
	s.inputChan = make(chan []envelope.Envelope, 10000)
	s.id, s.name = &id, conf.Name
	s.deliveryRequired, s.fanout = conf.DeliveryRequired, conf.Fanout
	return nil
}

func (s *Sink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
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

func (s *Sink) Enqueue(envelopes []envelope.Envelope) {
	s.inputChan <- envelopes
}

func (s *Sink) Shutdown() error {
	log.Debug().Msg("🟡 shutting down stdout sink")
	return nil
}
