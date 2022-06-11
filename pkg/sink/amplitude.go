package sink

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/request"
)

const AMPLITUDE_NA_ENDPOINT string = "https://api2.amplitude.com/2/httpapi"
const AMPLITUDE_EU_ENDPOINT string = "https://api.eu.amplitude.com/2/httpapi"

type amplitudeEvent struct {
	UserId             string                 `json:"user_id"`
	DeviceId           string                 `json:"device_id"`
	EventType          string                 `json:"event_type"`
	Time               int64                  `json:"time"`
	EventProperties    map[string]interface{} `json:"event_properties"`
	UserProperties     map[string]interface{} `json:"user_properties"`
	Groups             map[string]interface{} `json:"groups"`
	AppVersion         string                 `json:"app_version"`
	Platform           string                 `json:"platform"`
	OsName             string                 `json:"os_name"`
	OsVersion          string                 `json:"os_version"`
	DeviceBrand        string                 `json:"device_brand"`
	DeviceManufacturer string                 `json:"device_manufacturer"`
	DeviceModel        string                 `json:"device_model"`
	Carrier            string                 `json:"carrier"`
	Country            string                 `json:"country"`
	Region             string                 `json:"region"`
	City               string                 `json:"city"`
	Dma                string                 `json:"dma"`
	Language           string                 `json:"language"`
	Price              float32                `json:"price"`
	Quantity           int                    `json:"quantity"`
	Revenue            float32                `json:"revenue"`
	ProductId          string                 `json:"productId"`
	RevenueType        string                 `json:"revenueType"`
	LocationLat        float32                `json:"location_lat"`
	LocationLng        float32                `json:"location_lng"`
	Ip                 string                 `json:"ip"`
	Idfa               string                 `json:"idfa"`
	Idfv               string                 `json:"idfv"`
	Adid               string                 `json:"adid"`
	AndroidId          string                 `json:"android_id"`
	EventId            int                    `json:"event_id"`
	SessionId          uint32                 `json:"session_id"`
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
	var amplitudeEvents []amplitudeEvent
	for _, e := range envelopes {
		payload, _ := e.Payload.AsMap()
		evnt := amplitudeEvent{
			DeviceId:        *e.Device.Nid,
			EventType:       e.EventMeta.Namespace,
			Time:            e.Pipeline.Source.GeneratedTstamp.Unix(),
			EventProperties: payload,
			// UserProperties:  e.User,
			// Groups:
			// AppVersion: *e.Pipeline.Source.Version,
			// Platform:   *e.Pipeline.Source.Name, // FIXME
			// OsName:     *e.Device.Os.Name,
			// OsVersion:  *e.Device.Os.Version,
			// DeviceBrand: e.Device.,
			// DeviceManufacturer: ,
			// DeviceModel: ,
			// Carrier: ,
			// Country: ,
			// Region: ,
			// City: ,
			// Dma: ,
			// Language: ,
			// Price: ,
			// Quantity: ,
			// Revenue: ,
			// ProductId: ,
			// RevenueType: ,
			// LocationLat: ,
			// LocationLng: ,
			// Ip:        e.Device.Ip,
			// Idfa:      *e.Device.Idfa,
			// Idfv:      *e.Device.Idfv,
			// Adid:      *e.Device.AdId,
			// AndroidId: *e.Device.AndroidId,
			// EventId: ,
			// SessionId: ,
			InsertId: e.EventMeta.Uuid.String(),
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
	log.Debug().Msg("closing amplitude sink")
	// no-opo
}
