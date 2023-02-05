// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package manifold

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/annotator"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/privacy"
	"github.com/silverton-io/buz/pkg/registry"
	"github.com/silverton-io/buz/pkg/sink"
)

// A stupid-simple manifold with strict guarantees.
// This manifold requires buffering at the client level for substantial event volumes.
// Otherwise it will probably overload the configured sink(s).
type SimpleManifold struct {
	registry         *registry.Registry
	sinks            *[]sink.Sink
	conf             *config.Config
	collectorMetdata *meta.CollectorMeta
}

func (m *SimpleManifold) Initialize(registry *registry.Registry, sinks *[]sink.Sink, conf *config.Config, metadata *meta.CollectorMeta) error {
	m.registry = registry
	m.sinks = sinks
	m.conf = conf
	m.collectorMetdata = metadata
	return nil
}

func (m *SimpleManifold) Distribute(envelopes []envelope.Envelope) error {
	var validEnvelopes []envelope.Envelope
	var invalidEnvelopes []envelope.Envelope
	annotatedEnvelopes := annotator.Annotate(envelopes, m.registry)
	anonymizedEnvelopes := privacy.AnonymizeEnvelopes(annotatedEnvelopes, m.conf.Privacy)
	for _, e := range anonymizedEnvelopes {
		isValid := e.Validation.IsValid
		if *isValid {
			validEnvelopes = append(validEnvelopes, e)
		} else {
			invalidEnvelopes = append(invalidEnvelopes, e)
		}
	}

	for _, s := range *m.sinks {
		ctx := context.Background()
		if len(validEnvelopes) > 0 {
			log.Debug().Interface("sinkId", s.Id()).Interface("sinkName", s.Name()).Interface("deliveryRequired", s.DeliveryRequired()).Interface("sinkType", s.Type()).Msg("🟡 purging valid envelopes to sink")
			publishErr := s.BatchPublishValid(ctx, validEnvelopes)
			if publishErr != nil {
				log.Error().Err(publishErr).Interface("sinkId", s.Id()).Interface("sinkName", s.Name()).Interface("deliveryRequired", s.DeliveryRequired()).Interface("sinkType", s.Type()).Msg("🔴 could not purge valid envelopes to sink")
				if s.DeliveryRequired() {
					return publishErr
				}
			}
		}
		if len(invalidEnvelopes) > 0 {
			log.Debug().Interface("sinkId", s.Id()).Interface("sinkName", s.Name()).Interface("deliveryRequired", s.DeliveryRequired()).Interface("sinkType", s.Type()).Msg("🟡 purging invalid envelopes to sink")
			publishErr := s.BatchPublishInvalid(ctx, invalidEnvelopes)
			if publishErr != nil {
				log.Error().Err(publishErr).Interface("sinkId", s.Id()).Interface("sinkName", s.Name()).Interface("deliveryRequired", s.DeliveryRequired()).Interface("sinkType", s.Type()).Msg("🔴 could not purge invalid envelopes to sink")
				if s.DeliveryRequired() {
					return publishErr
				}
			}
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
