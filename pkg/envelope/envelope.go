// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"encoding/json"

	"github.com/silverton-io/buz/pkg/db"
)

const (
	VENDOR         string = "vendor"
	NAMESPACE      string = "namespace"
	VERSION        string = "version"
	FORMAT         string = "format"
	PATH           string = "path"
	INPUT_PROTOCOL string = "inputProtocol"
	SCHEMA         string = "schema"
)

type Envelope struct {
	db.BasePKeylessModel
	EventMeta    `json:"event" gorm:"type:json"`
	Pipeline     `json:"pipeline" gorm:"type:json"`
	Device       `json:"device,omitempty" gorm:"type:json"`
	*User        `json:"user,omitempty" gorm:"type:json"`
	*Session     `json:"session,omitempty" gorm:"type:json"`
	*Web         `json:"web,omitempty" gorm:"type:json"`
	*Annotations `json:"annotations,omitempty" gorm:"type:json"`
	*Enrichments `json:"enrichments,omitempty" gorm:"type:json"`
	Validation   `json:"validation" gorm:"type:json"`
	Contexts     *Contexts `json:"contexts,omitempty" gorm:"type:json"`
	Payload      Payload   `json:"payload" gorm:"type:json"`
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
	EventMeta    `json:"event" gorm:"type:jsonb"`
	Pipeline     `json:"pipeline" gorm:"type:jsonb"`
	*Device      `json:"device" gorm:"type:jsonb"`
	*User        `json:"user" gorm:"type:jsonb"`
	*Session     `json:"session" gorm:"type:jsonb"`
	*Web         `json:"web" gorm:"type:jsonb"`
	*Annotations `json:"annotations" gorm:"type:jsonb"`
	*Enrichments `json:"enrichments" gorm:"type:jsonb"`
	Validation   `json:"validation" gorm:"type:jsonb"`
	Contexts     *Contexts `json:"contexts" gorm:"type:jsonb"`
	Payload      Payload   `json:"payload" gorm:"type:jsonb"`
}

type StringEnvelope struct {
	db.BasePKeylessModel
	EventMeta    `json:"event" gorm:"type:string"`
	Pipeline     `json:"pipeline" gorm:"type:string"`
	*Device      `json:"device" gorm:"type:string"`
	*User        `json:"user" gorm:"type:string"`
	*Session     `json:"session" gorm:"type:string"`
	*Web         `json:"web" gorm:"type:string"`
	*Annotations `json:"annotations" gorm:"type:string"`
	*Enrichments `json:"enrichments" gorm:"type:string"`
	Validation   `json:"validation" gorm:"type:string"`
	Contexts     *Contexts `json:"contexts" gorm:"type:string"`
	Payload      Payload   `json:"payload" gorm:"type:string"`
}
