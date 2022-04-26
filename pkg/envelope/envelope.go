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
	SourceMetadata     `json:"sourceMetadata" gorm:"type:json"`
	CollectorMetadata  `json:"collectorMetadata" gorm:"type:json"`
	EventMetadata      `json:"eventMetadata" gorm:"type:json"`
	RelayMetadata      `json:"relayMetadata" gorm:"type:json"`
	ValidationMetadata `json:"validationMetadata" gorm:"type:json"`
	Payload            event.Event `json:"payload" gorm:"type:json"`
}

type PgEnvelope struct { // I really hate doing this - should find a better way to do dialect/db-specific types within the single envelope
	SourceMetadata     `json:"sourceMetadata" gorm:"type:jsonb"`
	CollectorMetadata  `json:"collectorMetadata" gorm:"type:jsonb"`
	EventMetadata      `json:"eventMetadata" gorm:"type:jsonb"`
	RelayMetadata      `json:"relayMetadata" gorm:"type:jsonb"`
	ValidationMetadata `json:"validationMetadata" gorm:"type:jsonb"`
	Payload            event.Event `json:"payload" gorm:"type:jsonb"`
}

type MysqlEnvelope struct { // I really hate doing this - should find a better way to do dialect/db-specific types within the single envelope
	SourceMetadata     `json:"sourceMetadata" gorm:"type:json"`
	CollectorMetadata  `json:"collectorMetadata" gorm:"type:json"`
	EventMetadata      `json:"eventMetadata" gorm:"type:json"`
	RelayMetadata      `json:"relayMetadata" gorm:"type:json"`
	ValidationMetadata `json:"validationMetadata" gorm:"type:json"`
	Payload            event.Event `json:"payload" gorm:"type:json"`
}

type ClickhouseEnvelope struct { // I really hate doing this - should find a better way to do dialect/db-specific types within the single envelope
	SourceMetadata     `json:"sourceMetadata" gorm:"type:string"`
	CollectorMetadata  `json:"collectorMetadata" gorm:"type:string"`
	EventMetadata      `json:"eventMetadata" gorm:"type:string"`
	RelayMetadata      `json:"relayMetadata" gorm:"type:string"`
	ValidationMetadata `json:"validationMetadata" gorm:"type:string"`
	Payload            event.Event `json:"payload" gorm:"type:string"`
}

type EventMetadata struct {
	Protocol           string    `json:"protocol"`
	Uuid               uuid.UUID `json:"uuid"`
	Vendor             string    `json:"vendor,omitempty"`
	PrimaryNamespace   string    `json:"primaryNamespace,omitempty"`
	SecondaryNamespace string    `json:"secondaryNamespace,omitempty"`
	TertiaryNamespace  string    `json:"tertiaryNamespace,omitempty"`
	Name               string    `json:"name,omitempty"`
	Version            string    `json:"version,omitempty"`
	Format             string    `json:"format,omitempty"`
	Path               string    `json:"path,omitempty"`
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
	Ip   string `json:"ip"`
}

type CollectorMetadata struct {
	Tstamp time.Time `json:"tstamp"`
}

type RelayMetadata struct {
	IsRelayed bool      `json:"isRelayed"`
	RelayId   uuid.UUID `json:"relayId"`
}

type ValidationMetadata struct {
	IsValid         bool             `json:"isValid"`
	ValidationError *ValidationError `json:"validationErrors"`
}
