package sink

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

const AMPLITUDE_NA_ENDPOINT string = "https://api2.amplitude.com/2/httpapi"
const AMPLITUDE_EU_ENDPOINT string = "https://api.eu.amplitude.com/2/httpapi"

type AmplitudeSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	endpoint         url.URL
	apiKey           string
}

func (s *AmplitudeSink) Id() *uuid.UUID {
	return s.id
}

func (s *AmplitudeSink) Name() string {
	return s.name
}

func (s *AmplitudeSink) Type() string {
	return AMPLITUDE
}

func (s *AmplitudeSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *AmplitudeSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing indicative sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	endpoint, err := url.Parse(AMPLITUDE_NA_ENDPOINT)
	if err != nil {
		return err
	}
	s.endpoint, s.apiKey = *endpoint, conf.AmplitudeApiKey
	return nil
}

func (s *AmplitudeSink) batchPublish(ctx context.Context, envelopes []envelope.Envelope) error {

	return nil
}

func (s *AmplitudeSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, envelopes)
	return err
}

func (s *AmplitudeSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, envelopes)
	return err
}

func (s *AmplitudeSink) Close() {
	log.Debug().Msg("closing amplitude sink")
	// no-opo
}
