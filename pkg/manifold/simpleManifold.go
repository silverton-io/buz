package manifold

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/sink"
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/silverton-io/honeypot/pkg/util"
)

// A stupid-simple manifold with strict guarantees
type SimpleManifold struct {
	sinks *[]sink.Sink
}

func (m *SimpleManifold) Initialize(sinks *[]sink.Sink) error {
	m.sinks = sinks
	return nil
}

func (m *SimpleManifold) Distribute(envelopes []envelope.Envelope, meta *tele.Meta) error {
	var validEnvelopes []envelope.Envelope
	var invalidEnvelopes []envelope.Envelope

	for _, e := range envelopes {
		util.Pprint(e)
		isValid := e.Validation.IsValid
		if isValid {
			meta.ProtocolStats.IncrementValid(&e.EventMeta, 1)
			validEnvelopes = append(validEnvelopes, e)
		} else {
			meta.ProtocolStats.IncrementInvalid(&e.EventMeta, 1)
			invalidEnvelopes = append(invalidEnvelopes, e)
		}
	}

	for _, s := range *m.sinks {
		ctx := context.Background()
		if len(validEnvelopes) > 0 {
			log.Debug().Interface("sinkId", s.Id()).Interface("sinkName", s.Name()).Interface("deliveryRequired", s.DeliveryRequired()).Interface("sinkType", s.Type()).Msg("purging valid envelopes to sink")
			publishErr := s.BatchPublishValid(ctx, validEnvelopes)
			if publishErr != nil {
				log.Error().Err(publishErr).Interface("sinkId", s.Id()).Interface("sinkName", s.Name()).Interface("deliveryRequired", s.DeliveryRequired()).Interface("sinkType", s.Type()).Msg("could not purge valid envelopes to sink")
				if s.DeliveryRequired() {
					return publishErr
				}
			}
		}
		if len(invalidEnvelopes) > 0 {
			log.Debug().Interface("sinkId", s.Id()).Interface("sinkName", s.Name()).Interface("deliveryRequired", s.DeliveryRequired()).Interface("sinkType", s.Type()).Msg("purging invalid envelopes to sink")
			publishErr := s.BatchPublishInvalid(ctx, invalidEnvelopes)
			if publishErr != nil {
				log.Error().Err(publishErr).Interface("sinkId", s.Id()).Interface("sinkName", s.Name()).Interface("deliveryRequired", s.DeliveryRequired()).Interface("sinkType", s.Type()).Msg("could not purge invalid envelopes to sink")
				if s.DeliveryRequired() {
					return publishErr
				}
			}
		}
	}
	return nil
}
