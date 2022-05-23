package envelope

import (
	"encoding/json"

	"github.com/silverton-io/honeypot/pkg/db"
	"github.com/silverton-io/honeypot/pkg/event"
)

const (
	VENDOR         string = "vendor"
	NAMESPACE      string = "namespace"
	VERSION        string = "version"
	FORMAT         string = "format"
	PATH           string = "path"
	INPUT_PROTOCOL string = "inputProtocol"
)

type Envelope struct {
	db.BasePKeylessModel
	EventMeta  `json:"eventMeta" gorm:"type:json"`
	Pipeline   `json:"pipeline" gorm:"type:json"`
	Device     `json:"device" gorm:"type:json"`
	User       `json:"user" gorm:"type:json"`
	Session    `json:"session" gorm:"type:json"`
	Page       `json:"page" gorm:"type:json"`
	Annotation `json:"annotation" gorm:"type:json"`
	Validation `json:"validation" gorm:"type:json"`
	Contexts   event.Contexts `json:"contexts" gorm:"type:json"`
	Payload    event.Event    `json:"payload" gorm:"type:json"`
}

func (e *Envelope) AsByte() ([]byte, error) {
	eBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return eBytes, nil
}

type JsonbEnvelope struct {
	db.BasePKeylessModel
	EventMeta  `json:"eventMeta" gorm:"type:jsonb"`
	Pipeline   `json:"pipeline" gorm:"type:jsonb"`
	Device     `json:"device" gorm:"type:jsonb"`
	User       `json:"user" gorm:"type:jsonb"`
	Session    `json:"session" gorm:"type:jsonb"`
	Page       `json:"page" gorm:"type:jsonb"`
	Annotation `json:"annotation" gorm:"type:jsonb"`
	Validation `json:"validation" gorm:"type:jsonb"`
	Contexts   event.Contexts `json:"contexts" gorm:"type:jsonb"`
	Payload    event.Event    `json:"payload" gorm:"type:jsonb"`
}

type StringEnvelope struct {
	db.BasePKeylessModel
	EventMeta  `json:"eventMeta" gorm:"type:string"`
	Pipeline   `json:"pipeline" gorm:"type:string"`
	Device     `json:"device" gorm:"type:string"`
	User       `json:"user" gorm:"type:string"`
	Session    `json:"session" gorm:"type:string"`
	Page       `json:"page" gorm:"type:string"`
	Annotation `json:"annotation" gorm:"type:string"`
	Validation `json:"validation" gorm:"type:string"`
	Contexts   event.Contexts `json:"contexts" gorm:"type:string"`
	Payload    event.Event    `json:"payload" gorm:"type:string"`
}
