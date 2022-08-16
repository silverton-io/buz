// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

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

const (
	AMPLITUDE_STANDARD_ENDPOINT string = "https://api2.amplitude.com/2/httpapi"
	AMPLITUDE_EU_ENDPOINT       string = "https://api.eu.amplitude.com/2/httpapi"
	AMPLITUDE_STANDARD          string = "standard"
	AMPLITUDE_EU                string = "eu"
)

type amplitudeEvent struct {
	UserId             string                 `json:"user_id"`
	DeviceId           string                 `json:"device_id"`
	EventType          string                 `json:"event_type"`
	Time               int64                  `json:"time"`
	EventProperties    map[string]interface{} `json:"event_properties,omitempty"`
	UserProperties     map[string]interface{} `json:"user_properties"`
	Groups             map[string]interface{} `json:"groups,omitempty"`
	AppVersion         *string                `json:"app_version,omitempty"`
	Platform           *string                `json:"platform,omitempty"`
	OsName             *string                `json:"os_name,omitempty"`
	OsVersion          *string                `json:"os_version,omitempty"`
	DeviceBrand        *string                `json:"device_brand,omitempty"`
	DeviceManufacturer *string                `json:"device_manufacturer,omitempty"`
	DeviceModel        *string                `json:"device_model,omitempty"`
	Carrier            *string                `json:"carrier,omitempty"`
	Country            *string                `json:"country,omitempty"`
	Region             *string                `json:"region,omitempty"`
	City               *string                `json:"city,omitempty"`
	Dma                *string                `json:"dma,omitempty"`
	Language           *string                `json:"language,omitempty"`
	Price              *float32               `json:"price,omitempty"`
	Quantity           *int                   `json:"quantity,omitempty"`
	Revenue            *float32               `json:"revenue,omitempty"`
	ProductId          *string                `json:"productId,omitempty"`
	RevenueType        *string                `json:"revenueType,omitempty"`
	LocationLat        *float32               `json:"location_lat,omitempty"`
	LocationLng        *float32               `json:"location_lng,omitempty"`
	Ip                 string                 `json:"ip,omitempty"`
	Idfa               *string                `json:"idfa,omitempty"`
	Idfv               *string                `json:"idfv,omitempty"`
	Adid               *string                `json:"adid,omitempty"`
	AndroidId          *string                `json:"android_id,omitempty"`
	EventId            int                    `json:"event_id,omitempty"`
	SessionId          *uint32                `json:"session_id,omitempty"`
	InsertId           string                 `json:"insert_id"`
	// TODO
	// plan
	// plan.branch
	// plan.source
	// plan.version
}

type amplitudeEventBatch struct {
	ApiKey string           `json:"api_key"`
	Events []amplitudeEvent `json:"events"`
}

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
	log.Debug().Msg("🟡 initializing indicative sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	var e string
	if conf.AmplitudeRegion == AMPLITUDE_EU {
		e = AMPLITUDE_EU_ENDPOINT
	} else {
		e = AMPLITUDE_STANDARD_ENDPOINT
	}
	endpoint, err := url.Parse(e)
	if err != nil {
		return err
	}
	s.endpoint, s.apiKey = *endpoint, conf.AmplitudeApiKey
	return nil
}

func (s *AmplitudeSink) batchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	var amplitudeEvents []amplitudeEvent
	for _, e := range envelopes {
		mappedEnvelope, _ := e.AsMap()
		flattenedEnvelope, err := flatten.Flatten(mappedEnvelope, "", flatten.DotStyle)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not flatten payload")
			return err
		}
		evnt := amplitudeEvent{
			Time:            e.Pipeline.Source.GeneratedTstamp.Unix(),
			DeviceId:        e.Device.Id,
			EventType:       e.EventMeta.Namespace,
			EventProperties: flattenedEnvelope,
			InsertId:        e.EventMeta.Uuid.String(),
		}
		if e.User.Id != nil {
			evnt.UserId = *e.User.Id
		}
		amplitudeEvents = append(amplitudeEvents, evnt)
	}
	payload := amplitudeEventBatch{
		ApiKey: s.apiKey,
		Events: amplitudeEvents,
	}
	_, err := request.PostPayload(s.endpoint, payload) // FIXME! Do something better with non-200 responses
	if err != nil {
		return err
	}
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
	log.Debug().Msg("🟡 closing amplitude sink")
	// no-opo
}
