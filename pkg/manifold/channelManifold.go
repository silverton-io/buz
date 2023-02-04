// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package manifold

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/annotator"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/params"
	"github.com/silverton-io/buz/pkg/privacy"
	"github.com/silverton-io/buz/pkg/registry"
	"github.com/silverton-io/buz/pkg/sink"
	"github.com/silverton-io/buz/pkg/util"
)

type ChannelManifold struct {
	registry      *registry.Registry
	sinks         *[]sink.Sink
	handlerParams *params.Handler
	invalid       chan envelope.Envelope
	valid         chan envelope.Envelope
	shutdown      chan int
}

func (m *ChannelManifold) Initialize(registry *registry.Registry, sinks *[]sink.Sink, handlerParams *params.Handler) error {
	m.registry = registry
	m.sinks = sinks
	m.handlerParams = handlerParams
	m.valid = make(chan envelope.Envelope)
	m.invalid = make(chan envelope.Envelope)
	m.shutdown = make(chan int, 1)
	log.Debug().Msg("spinning up manifold goroutine")
	go func(e <-chan envelope.Envelope, quit chan int) {
		for {
			select {
			case envelope := <-e:
				util.Pprint(envelope)
			case <-quit:
				log.Info().Msg("manifold shut down")
				return
			}
		}
	}(m.valid, m.shutdown)
	return nil
}

func (m *ChannelManifold) Distribute(envelopes []envelope.Envelope) error {
	annotatedEnvelopes := annotator.Annotate(envelopes, m.registry)
	anonymizedEnvelopes := privacy.AnonymizeEnvelopes(annotatedEnvelopes, m.handlerParams.Config.Privacy)
	for _, e := range anonymizedEnvelopes {
		isValid := e.Validation.IsValid
		log.Debug().Interface("payload", e).Msg("sending msg to chan")
		if *isValid {
			m.valid <- e
		} else {
			m.invalid <- e
		}
	}
	return nil
}

func (m *ChannelManifold) Shutdown() error {
	log.Info().Msg("shutting down channel manifold")
	m.shutdown <- 1
	return nil
}
