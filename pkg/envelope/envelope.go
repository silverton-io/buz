// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/config"
)

const (
	PROTOCOL  string = "protocol"
	VENDOR    string = "vendor"
	NAMESPACE string = "namespace"
	VERSION   string = "version"
	SCHEMA    string = "schema"
	IS_VALID  string = "isValid"
)

// An envelope consisting of minimally-defined properties
type Envelope struct {
	Uuid            uuid.UUID        `json:"uuid"`
	Timestamp       time.Time        `json:"timestamp" sql:"index"`
	BuzTimestamp    time.Time        `json:"buzTimestamp" sql:"index"`
	BuzVersion      string           `json:"buzVersion"`
	BuzName         string           `json:"buzName"`
	BuzEnv          string           `json:"buzEnv"`
	Protocol        string           `json:"protocol"`
	Schema          string           `json:"schema"`
	Vendor          string           `json:"vendor"`
	Namespace       string           `json:"namespace"`
	Version         string           `json:"version"`
	IsValid         bool             `json:"isValid"`
	ValidationError *ValidationError `json:"validationError,omitempty" gorm:"type:json"`
	Contexts        *Contexts        `json:"contexts,omitempty" gorm:"type:json"`
	Payload         Payload          `json:"payload" gorm:"type:json"`
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
	Uuid            uuid.UUID        `json:"uuid" gorm:"type:uuid"`
	Timestamp       time.Time        `json:"timestamp" sql:"index"`
	BuzTimestamp    time.Time        `json:"buzTimestamp" sql:"index"`
	BuzVersion      string           `json:"buzVersion"`
	BuzName         string           `json:"buzName"`
	BuzEnv          string           `json:"buzEnv"`
	Protocol        string           `json:"protocol"`
	Schema          string           `json:"schema"`
	Vendor          string           `json:"vendor"`
	Namespace       string           `json:"namespace"`
	Version         string           `json:"version"`
	IsValid         bool             `json:"isValid"`
	ValidationError *ValidationError `json:"validationError,omitempty" gorm:"type:jsonb"`
	Contexts        *Contexts        `json:"contexts,omitempty" gorm:"type:jsonb"`
	Payload         Payload          `json:"payload" gorm:"type:jsonb"`
}

type StringEnvelope struct {
	Uuid            uuid.UUID        `json:"uuid" gorm:"type:uuid"`
	Timestamp       time.Time        `json:"timestamp" sql:"index"`
	BuzTimestamp    time.Time        `json:"buzTimestamp" sql:"index"`
	BuzVersion      string           `json:"buzVersion"`
	BuzName         string           `json:"buzName"`
	BuzEnv          string           `json:"buzEnv"`
	Protocol        string           `json:"protocol"`
	Schema          string           `json:"schema"`
	Vendor          string           `json:"vendor"`
	Namespace       string           `json:"namespace"`
	Version         string           `json:"version"`
	IsValid         bool             `json:"isValid"`
	ValidationError *ValidationError `json:"validationError,omitempty" gorm:"type:string"`
	Contexts        *Contexts        `json:"contexts,omitempty" gorm:"type:string"`
	Payload         Payload          `json:"payload" gorm:"type:string"`
}

// Build a new envelope with base fields populated
func NewEnvelope(conf config.App) Envelope {
	now := time.Now().UTC()
	envelope := Envelope{
		Uuid:         uuid.New(),
		Timestamp:    now,
		BuzTimestamp: now,
		BuzVersion:   conf.Version,
		BuzName:      conf.Name,
		BuzEnv:       conf.Env,
	}
	return envelope
}
