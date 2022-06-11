package sink

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/jeremywohl/flatten/v2"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/request"
)

const INDICATIVE_BATCH_ENDPOINT string = "https://api.indicative.com/service/event/batch"

type indicativeEvent struct {
	EventName     string                 `json:"eventName"`
	EventUniqueId string                 `json:"eventUniqueId"`
	Properties    map[string]interface{} `json:"properties"`
	EventTime     int64                  `json:"eventTime"`
}

type indicativeEventBatch struct {
	ApiKey string            `json:"apiKey"`
	Events []indicativeEvent `json:"events"`
}

type IndicativeSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	endpoint         url.URL
	apiKey           string
}

func (s *IndicativeSink) Id() *uuid.UUID {
	return s.id
}

func (s *IndicativeSink) Name() string {
	return s.name
}

func (s *IndicativeSink) Type() string {
	return INDICATIVE
}

func (s *IndicativeSink) DeliveryRequired() bool {
	return s.deliveryRequired
}

func (s *IndicativeSink) Initialize(conf config.Sink) error {
	log.Debug().Msg("initializing indicative sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	endpoint, err := url.Parse(INDICATIVE_BATCH_ENDPOINT)
	if err != nil {
		return err
	}
	s.endpoint, s.apiKey = *endpoint, conf.IndicativeApiKey
	return nil
}

func (s *IndicativeSink) batchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	var indicativeEvents []indicativeEvent
	for _, e := range envelopes {
		propertyMap, err := e.AsMap()
		if err != nil {
			log.Error().Err(err).Msg("could not coerce envelope to map")
			return err
		}
		flattenedPropertyMap, err := flatten.Flatten(propertyMap, "", flatten.UnderscoreStyle) // NOTE! Indicative does not allow nested properties.
		if err != nil {
			log.Error().Err(err).Msg("could not flatten properties")
		}
		evnt := indicativeEvent{
			EventName:     e.EventMeta.Namespace,
			EventUniqueId: e.EventMeta.Uuid.String(),
			Properties:    flattenedPropertyMap,
			EventTime:     e.Source.GeneratedTstamp.Unix(),
		}
		indicativeEvents = append(indicativeEvents, evnt)
	}
	payload := indicativeEventBatch{
		ApiKey: s.apiKey,
		Events: indicativeEvents,
	}
	_, err := request.PostPayload(s.endpoint, payload)
	if err != nil {
		return err
	}
	return nil
}

func (s *IndicativeSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, envelopes)
	return err
}

func (s *IndicativeSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) error {
	err := s.batchPublish(ctx, envelopes)
	return err
}

func (s *IndicativeSink) Close() {
	log.Debug().Msg("closing indicative sink")
	// no-opo
}
