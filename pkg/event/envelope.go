package event

import (
	"time"

	"github.com/silverton-io/honeypot/pkg/validator"
)

type Envelope struct {
	EventProtocol    string                       `json:"eventProtocol"`
	Tstamp           time.Time                    `json:"tstamp"`
	Ip               string                       `json:"ip"`
	IsValid          *bool                        `json:"isValid"`
	ValidationErrors *[]validator.ValidationError `json:"validationErrors"`
	Event            map[string]interface{}       `json:"event"`
}
