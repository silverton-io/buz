package envelope

import (
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/event"
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

type Envelope struct {
	Id              uuid.UUID      `json:"id"`
	EventProtocol   string         `json:"eventProtocol"`
	EventMetadata   *EventMetadata `json:"eventMetadata"`
	SourceMetadata  `json:"sourceMetadata"`
	Tstamp          time.Time        `json:"tstamp"`
	Ip              string           `json:"ip"`
	IsValid         *bool            `json:"isValid"`
	IsRelayed       *bool            `json:"isRelayed"`
	RelayedId       *uuid.UUID       `json:"relayedId"`
	ValidationError *ValidationError `json:"validationErrors"`
	Payload         event.Event      `json:"payload"`
}

type EventMetadata struct {
	Vendor            string `json:"vendor,omitempty"`
	PrimaryCategory   string `json:"primaryCategory,omitempty"`
	SecondaryCategory string `json:"secondaryCategory,omitempty"`
	TertiaryCategory  string `json:"tertiaryCategory,omitempty"`
	Name              string `json:"name,omitempty"`
	Version           string `json:"version,omitempty"`
	Format            string `json:"format,omitempty"`
	Path              string `json:"path,omitempty"`
}

type SourceMetadata struct {
	Name string `json:"name,omitempty"`
}
