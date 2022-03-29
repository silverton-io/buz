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
	Id              uuid.UUID        `json:"id"`
	EventProtocol   string           `json:"eventProtocol"`
	EventSchema     string           `json:"eventSchema"`
	Source          string           `json:"source"`
	Tstamp          time.Time        `json:"tstamp"`
	Ip              string           `json:"ip"`
	IsValid         *bool            `json:"isValid"`
	IsRelayed       *bool            `json:"isRelayed"`
	RelayedId       *uuid.UUID       `json:"relayedId"`
	ValidationError *ValidationError `json:"validationErrors"`
	Payload         event.Event      `json:"payload"`
}
