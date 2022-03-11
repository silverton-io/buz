package event

import (
	"time"

	"github.com/silverton-io/honeypot/pkg/validator"
)

type HoneypotEnvelopeMetadata struct {
	EventProtocol     string                    `json:"eventProtocol"`
	EventVendor       string                    `json:"eventVendor"`
	EventId           string                    `json:"eventId"`
	EventVersion      string                    `json:"eventVersion"`
	SchemaFormat      string                    `json:"schemaFormat"`
	IsValid           bool                      `json:"isValid"`
	ValidationError   validator.ValidationError `json:"validationError"`
	HoneypotTimestamp time.Time                 `json:"honeypotTimestamp"`
}

type HoneypotEnvelope struct {
	HoneypotEnvelopeMetadata `json:"hMetadata"`
	HPayload                 map[string]interface{} `json:"hPayload"`
}

type SelfDescribingEnvelope struct {
	Contexts []SelfDescribingContext `json:"contexts"`
	Payload  SelfDescribingPayload   `json:"payload"`
}

type SelfDescribingPayload struct {
	Schema string                 `json:"schema"`
	Data   map[string]interface{} `json:"data"`
}

type SelfDescribingContext SelfDescribingPayload
