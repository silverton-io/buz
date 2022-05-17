package envelope

import (
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
	INPUT_PROTOCOL            string = "inputProtocol"
)

type Envelope struct {
	Event      `json:"event" gorm:"type:json"`
	Pipeline   `json:"pipeline" gorm:"type:json"`
	Device     `json:"device" gorm:"type:json"`
	User       `json:"user" gorm:"type:json"`
	Session    `json:"session" gorm:"type:json"`
	Page       `json:"page" gorm:"type:json"`
	Validation `json:"validation" gorm:"type:json"`
	Contexts   event.Contexts `json:"contexts" gorm:"type:json"`
	Payload    event.Event    `json:"payload" gorm:"type:json"`
}

type JsonbEnvelope struct {
	Event      `json:"event" gorm:"type:jsonb"`
	Pipeline   `json:"pipeline" gorm:"type:jsonb"`
	Device     `json:"device" gorm:"type:jsonb"`
	User       `json:"user" gorm:"type:jsonb"`
	Session    `json:"session" gorm:"type:jsonb"`
	Page       `json:"page" gorm:"type:jsonb"`
	Validation `json:"validation" gorm:"type:jsonb"`
	Contexts   event.Contexts `json:"contexts" gorm:"type:jsonb"`
	Payload    event.Event    `json:"payload" gorm:"type:jsonb"`
}

type StringEnvelope struct {
	Event      `json:"event" gorm:"type:string"`
	Pipeline   `json:"pipeline" gorm:"type:string"`
	Device     `json:"device" gorm:"type:string"`
	User       `json:"user" gorm:"type:string"`
	Session    `json:"session" gorm:"type:string"`
	Page       `json:"page" gorm:"type:string"`
	Validation `json:"validation" gorm:"type:string"`
	Contexts   event.Contexts `json:"contexts" gorm:"type:string"`
	Payload    event.Event    `json:"payload" gorm:"type:string"`
}
