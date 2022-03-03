package event

import "github.com/silverton-io/honeypot/pkg/validator"

type InvalidEvent struct {
	ValidationError *validator.ValidationError `json:"validationError"`
	Event           interface{}                `json:"event"`
}
