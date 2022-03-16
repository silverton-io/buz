package event

import (
	"time"
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
	EventProtocol   string           `json:"eventProtocol"`
	EventSchema     *string          `json:"eventSchema"`
	Source          string           `json:"source"`
	Tstamp          time.Time        `json:"tstamp"`
	Ip              string           `json:"ip"`
	IsValid         *bool            `json:"isValid"`
	ValidationError *ValidationError `json:"validationErrors"`
	Payload         Event            `json:"payload"`
}
