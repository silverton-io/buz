// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"encoding/json"

	"github.com/silverton-io/buz/pkg/db"
	"github.com/silverton-io/buz/pkg/event"
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
	Web        `json:"web" gorm:"type:json"`
	Annotation `json:"annotation" gorm:"type:json"`
	Validation `json:"validation" gorm:"type:json"`
	Contexts   event.Contexts `json:"contexts" gorm:"type:json"`
	Payload    event.Payload  `json:"payload" gorm:"type:json"`
}

func (e *Envelope) AsMap() (map[string]interface{}, error) {
	var m map[string]interface{}
	marshaledEnvelope, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(marshaledEnvelope, &m); err != nil {
		return nil, err
	}
	return m, nil
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
	Web        `json:"web" gorm:"type:jsonb"`
	Annotation `json:"annotation" gorm:"type:jsonb"`
	Validation `json:"validation" gorm:"type:jsonb"`
	Contexts   event.Contexts `json:"contexts" gorm:"type:jsonb"`
	Payload    event.Payload  `json:"payload" gorm:"type:jsonb"`
}

type StringEnvelope struct {
	db.BasePKeylessModel
	EventMeta  `json:"eventMeta" gorm:"type:string"`
	Pipeline   `json:"pipeline" gorm:"type:string"`
	Device     `json:"device" gorm:"type:string"`
	User       `json:"user" gorm:"type:string"`
	Session    `json:"session" gorm:"type:string"`
	Web        `json:"web" gorm:"type:string"`
	Annotation `json:"annotation" gorm:"type:string"`
	Validation `json:"validation" gorm:"type:string"`
	Contexts   event.Contexts `json:"contexts" gorm:"type:string"`
	Payload    event.Payload  `json:"payload" gorm:"type:string"`
}
