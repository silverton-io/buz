package sink

import (
	"context"
	"net/url"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
)

const (
	AMPLITUDE_SERVER_STANDARD = "standard"
	AMPLITUDE_SERVER_EU       = "eu"
)

type AmplitudeEvent struct {
	UserId             string                 `json:"user_id"`
	DeviceId           string                 `json:"device_id"`
	EventType          string                 `json:"event_type"`
	Time               uint64                 `json:"time,omitempty"`
	EventProperties    map[string]interface{} `json:"event_properties,omitempty"`
	UserProperties     map[string]interface{} `json:"user_properties,omitempty"`
	Groups             map[string]interface{} `json:"groups,omitempty"`
	AppVersion         string                 `json:"app_version,omitempty"`
	Platform           string                 `json:"platform,omitempty"`
	OsName             string                 `json:"os_name,omitempty"`
	OsVersion          string                 `json:"os_version,omitempty"`
	DeviceBrand        string                 `json:"device_brand,omitempty"`
	DeviceManufacturer string                 `json:"device_manufacturer,omitempty"`
	DeviceModel        string                 `json:"device_model,omitempty"`
	Carrier            string                 `json:"carrier,omitempty"`
	Country            string                 `json:"country,omitempty"`
	Region             string                 `json:"region,omitempty"`
	City               string                 `json:"city,omitempty"`
	Dma                string                 `json:"dma,omitempty"`
	Language           string                 `json:"language,omitempty"`
	Price              float32                `json:"price,omitempty"`
	Quantity           int                    `json:"quantity,omitempty"`
	Revenue            float32                `json:"revenue,omitempty"`
	ProductId          string                 `json:"productId,omitempty"`
	RevenueType        string                 `json:"revenueType,omitempty"`
	LocationLat        float32                `json:"location_lat,omitempty"`
	LocationLng        float32                `json:"location_lng,omitempty"`
	Ip                 string                 `json:"ip,omitempty"`
	Idfa               string                 `json:"idfa,omitempty"`
	Idfv               string                 `json:"idfv,omitempty"`
	Adid               string                 `json:"adid,omitempty"`
	AndroidId          string                 `json:"android_id,omitempty"`
	EventId            int                    `json:"event_id,omitempty"`
	SessionId          uint64                 `json:"session_id,omitempty"`
	InsertId           string                 `json:"insert_id,omitempty"`
	Plan               map[string]interface{} `json:"plan,omitempty"`
}

type AmplitudePayload struct {
	ApiKey  string                 `json:"api_key"`
	Events  []AmplitudeEvent       `json:"events"`
	Options map[string]interface{} `json:"options,omitempty"`
}

func buildAmplitudeEndpoint(server string) (*url.URL, error) {
	switch server {
	case AMPLITUDE_SERVER_EU:
		url, err := url.Parse("https://api.eu.amplitude.com/2/httpapi")
		return url, err
	default:
		url, err := url.Parse("https://api2.amplitude.com/2/httpapi")
		return url, err
	}
}

func buildAmplitudePayloadFromEnvelopes(apiKey string, envelopes []envelope.Envelope) AmplitudePayload {
	var events []AmplitudeEvent
	for _, env := range envelopes {
		payload, err := env.Payload.AsMap()
		if err != nil {
			log.Error().Err(err).Msg("could not coerce envelope payload to map")
		}
		event := AmplitudeEvent{
			UserId:          *env.UserMetadata.Uid,
			DeviceId:        *env.UserMetadata.Duid,
			EventType:       env.EventMetadata.Name,
			EventProperties: payload,
		}
		events = append(events, event)
	}
	pl := AmplitudePayload{
		ApiKey: apiKey,
		Events: events,
	}
	return pl
}

type AmplitudeSink struct {
	id               *uuid.UUID
	name             string
	deliveryRequired bool
	apiKey           string
	apiEndpoint      url.URL
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
	log.Debug().Msg("initializing amplitude sink")
	id := uuid.New()
	endpoint, err := buildAmplitudeEndpoint(conf.AmplitudeServer)
	if err != nil {
		log.Error().Err(err).Msg("could not build amplitude api endpoint")
		return err
	}
	s.id, s.name, s.deliveryRequired = &id, conf.Name, conf.DeliveryRequired
	s.apiKey, s.apiEndpoint = conf.AmplitudeApiKey, *endpoint
	return nil
}

func (s *AmplitudeSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) error {
	return nil
}

func (s *AmplitudeSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) error {
	return nil
}

func (s *AmplitudeSink) Close() {
	log.Debug().Msg("closing amplitude sink")
}
