// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package manifold

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/annotator"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/privacy"
	"github.com/silverton-io/buz/pkg/registry"
	"github.com/silverton-io/buz/pkg/sink"
	"github.com/silverton-io/buz/pkg/util"
)

type ChannelManifold struct {
	registry      *registry.Registry
	sinks         *[]sink.Sink
	conf          *config.Config
	collectorMeta *meta.CollectorMeta
	inputChan     chan []envelope.Envelope
	shutdownChan  chan int
}

func (m *ChannelManifold) Initialize(registry *registry.Registry, sinks *[]sink.Sink, conf *config.Config, metadata *meta.CollectorMeta) error {
	m.registry = registry
	m.sinks = sinks
	m.conf = conf
	m.collectorMeta = metadata
	m.inputChan = make(chan []envelope.Envelope)
	m.shutdownChan = make(chan int, 1)
	log.Debug().Msg("spinning up manifold goroutine")
	go func(e <-chan []envelope.Envelope, quit chan int) {
		for {
			select {
			case envelopes := <-e:
				// FIXME!!! Make this actually do something
				util.Pprint(envelopes)
			case <-quit:
				log.Info().Msg("manifold shut down")
				return
			}
		}
	}(m.inputChan, m.shutdownChan)
	return nil
}

func (m *ChannelManifold) Distribute(envelopes []envelope.Envelope) error {
	annotatedEnvelopes := annotator.Annotate(envelopes, m.registry)
	anonymizedEnvelopes := privacy.AnonymizeEnvelopes(annotatedEnvelopes, m.conf.Privacy)
	log.Debug().Interface("payload", anonymizedEnvelopes).Msg("sending envelopes to chan")
	m.inputChan <- anonymizedEnvelopes
	return nil
}

func (m *ChannelManifold) GetRegistry() *registry.Registry {
	return m.registry
}

func (m *ChannelManifold) Shutdown() error {
	log.Info().Msg("shutting down channel manifold")
	m.shutdownChan <- 1
	return nil
}
