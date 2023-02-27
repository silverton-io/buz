// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package amplitude

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/jeremywohl/flatten/v2"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/request"
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

type Sink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	endpoint         url.URL
	apiKey           string
	inputChan        chan []envelope.Envelope
	shutdownChan     chan int
}

func (s *Sink) Metadata() backendutils.SinkMetadata {
	sinkType := "amplitude"
	return backendutils.SinkMetadata{
		Id:               s.id,
		Name:             s.name,
		Type:             sinkType,
		DeliveryRequired: s.deliveryRequired,
	}
}

func (s *Sink) Initialize(conf config.Sink) error {
	log.Debug().Msg("🟡 initializing indicative sink")
	id := uuid.New()
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.inputChan = make(chan []envelope.Envelope, 10000)
	s.shutdownChan = make(chan int, 1)
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
	go func(s *Sink) {
		for {
			select {
			case envelopes := <-s.inputChan:
				ctx := context.Background()
				s.Dequeue(ctx, envelopes)
			case <-s.shutdownChan:
				err := s.Shutdown()
				if err != nil {
					log.Error().Err(err).Interface("metadata", s.Metadata()).Msg("sink did not safely shut down")
				}
				return
			}
		}
	}(s)
	return nil
}

func (s *Sink) batchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
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

func (s *Sink) Enqueue(envelopes []envelope.Envelope) {
	log.Debug().Interface("metadata", s.Metadata()).Msg("enqueueing envelopes")
	s.inputChan <- envelopes
}

func (s *Sink) Dequeue(ctx context.Context, envelopes []envelope.Envelope) error {
	log.Debug().Interface("metadata", s.Metadata()).Msg("dequeueing envelopes")
	err := s.batchPublish(ctx, envelopes)
	return err
}

func (s *Sink) Shutdown() error {
	log.Debug().Msg("🟢 shutting down amplitude sink")
	s.shutdownChan <- 1
	return nil
}
