package envelope

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/event"
)

const (
	EVENT_VENDOR              string = "vendor"
	EVENT_PRIMARY_NAMESPACE   string = "primaryNamespace"
	EVENT_SECONDARY_NAMESPACE string = "secondaryNamespace"
	EVENT_TERTIARY_NAMESPACE  string = "tertiaryNamespace"
	EVENT_NAME                string = "name"
	EVENT_VERSION             string = "version"
	EVENT_FORMAT              string = "format"
	EVENT_PATH                string = "path"
)

type PayloadValidationError struct {
	Field       string `json:"field"`
	Description string `json:"description"`
	ErrorType   string `json:"errorType"`
}

type ValidationError struct {
	ErrorType       *string                  `json:"errorType"`
	ErrorResolution *string                  `json:"errorResolution"`
	Errors          []PayloadValidationError `json:"payloadValidationErrors"`
}

func (e *ValidationError) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e *ValidationError) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), e)
}

type Envelope struct {
	// gorm.Model FIXME: maybe at some point in the future?
	Id              uuid.UUID      `json:"id"`
	EventProtocol   string         `json:"eventProtocol"`
	EventMetadata   *EventMetadata `json:"eventMetadata" gorm:"type:json"`
	SourceMetadata  `json:"sourceMetadata" gorm:"type:json"`
	Tstamp          time.Time        `json:"tstamp"`
	Ip              string           `json:"ip"`
	IsValid         *bool            `json:"isValid"`
	IsRelayed       *bool            `json:"isRelayed"`
	RelayedId       *uuid.UUID       `json:"relayedId"`
	ValidationError *ValidationError `json:"validationErrors" gorm:"type:json"`
	Payload         event.Event      `json:"payload" gorm:"type:json"`
}

type EventMetadata struct {
	Vendor             string `json:"vendor,omitempty"`
	PrimaryNamespace   string `json:"primaryNamespace,omitempty"`
	SecondaryNamespace string `json:"secondaryNamespace,omitempty"`
	TertiaryNamespace  string `json:"tertiaryNamespace,omitempty"`
	Name               string `json:"name,omitempty"`
	Version            string `json:"version,omitempty"`
	Format             string `json:"format,omitempty"`
	Path               string `json:"path,omitempty"`
}

func (e *EventMetadata) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e *EventMetadata) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), e)
}

type SourceMetadata struct {
	Name string `json:"name,omitempty"`
}
