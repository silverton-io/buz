// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
	Uuid               uuid.UUID        `json:"uuid"`
	Timestamp          time.Time        `json:"timestamp" sql:"index"`
	CollectorTimestamp time.Time        `json:"collectorTimestamp" sql:"index"`
	Protocol           string           `json:"protocol"`
	Schema             string           `json:"schema"`
	Vendor             string           `json:"vendor"`
	Namespace          string           `json:"namespace"`
	Version            string           `json:"version"`
	IsValid            bool             `json:"isValid"`
	ValidationError    *ValidationError `json:"validationError,omitempty"`
	Contexts           *Contexts        `json:"contexts,omitempty"`
	Payload            Payload          `json:"payload" gorm:"type:json"`
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
	Uuid               uuid.UUID        `json:"uuid"`
	Timestamp          time.Time        `json:"timestamp" sql:"index"`
	CollectorTimestamp time.Time        `json:"collectorTimestamp" sql:"index"`
	Protocol           string           `json:"protocol"`
	Schema             string           `json:"schema"`
	Vendor             string           `json:"vendor"`
	Namespace          string           `json:"namespace"`
	Version            string           `json:"version"`
	IsValid            bool             `json:"isValid"`
	ValidationError    *ValidationError `json:"validationError,omitempty"`
	Contexts           *Contexts        `json:"contexts,omitempty"`
	Payload            Payload          `json:"payload" gorm:"type:jsonb"`
}

type StringEnvelope struct {
	Uuid               uuid.UUID        `json:"uuid"`
	Timestamp          time.Time        `json:"timestamp" sql:"index"`
	CollectorTimestamp time.Time        `json:"collectorTimestamp" sql:"index"`
	Protocol           string           `json:"protocol"`
	Schema             string           `json:"schema"`
	Vendor             string           `json:"vendor"`
	Namespace          string           `json:"namespace"`
	Version            string           `json:"version"`
	IsValid            bool             `json:"isValid"`
	ValidationError    *ValidationError `json:"validationError,omitempty"`
	Contexts           *Contexts        `json:"contexts,omitempty"`
	Payload            Payload          `json:"payload" gorm:"type:string"`
}

// Build a new envelope with base fields populated
func NewEnvelope() Envelope {
	now := time.Now().UTC()
	envelope := Envelope{
		Uuid:               uuid.New(),
		Timestamp:          now,
		CollectorTimestamp: now,
	}
	return envelope
}
