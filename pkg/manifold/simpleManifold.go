// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package manifold

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/annotator"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/privacy"
	"github.com/silverton-io/buz/pkg/registry"
)

// A stupid-simple manifold with strict guarantees.
// This manifold requires buffering at the client level for substantial event volumes.
// Otherwise it will probably overload the configured sink(s).
type SimpleManifold struct {
	registry         *registry.Registry
	sinks            *[]backendutils.Sink
	conf             *config.Config
	collectorMetdata *meta.CollectorMeta
}

func (m *SimpleManifold) Initialize(registry *registry.Registry, sinks *[]backendutils.Sink, conf *config.Config, metadata *meta.CollectorMeta) error {
	m.registry = registry
	m.sinks = sinks
	m.conf = conf
	m.collectorMetdata = metadata
	return nil
}

func (m *SimpleManifold) Enqueue(envelopes []envelope.Envelope) error {
	annotatedEnvelopes := annotator.Annotate(envelopes, m.registry)
	anonymizedEnvelopes := privacy.AnonymizeEnvelopes(annotatedEnvelopes, m.conf.Privacy)
	for _, sink := range *m.sinks {
		meta := sink.Metadata()
		log.Debug().Interface("metadata", meta).Msg("🟡 enqueueing envelopes to sink")
		err := sink.Enqueue(anonymizedEnvelopes)
		if err != nil {
			log.Error().Err(err).Interface("metadata", sink.Metadata()).Msg("failed to enqueue envelopes to sink")
		}
	}
	return nil
}

func (m *SimpleManifold) GetRegistry() *registry.Registry {
	return m.registry
}

func (m *SimpleManifold) Shutdown() error {
	log.Info().Msg("shutting down simple manifold")
	log.Info().Msg("manifold shut down")
	return nil
}
