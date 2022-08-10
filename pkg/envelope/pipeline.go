// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package envelope

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Pipeline struct {
	Source    `json:"source,omitempty"`
	Collector `json:"collector,omitempty"`
	Relay     `json:"relay,omitempty"`
}

func (p Pipeline) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p Pipeline) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &p)
}

type Source struct {
	GeneratedTstamp time.Time `json:"generatedTstamp,omitempty"`
	SentTstamp      time.Time `json:"sentTstamp,omitempty"`
	Name            *string   `json:"name,omitempty"`
	Version         *string   `json:"version,omitempty"`
}

func (s Source) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

func (s Source) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &s)
}

type Collector struct {
	Tstamp  time.Time `json:"tstamp,omitempty"`
	Name    *string   `json:"name,omitempty"`
	Version *string   `json:"version,omitempty"`
}

func (c Collector) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c Collector) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &c)
}

type Relay struct {
	Relayed bool       `json:"relayed,omitempty"`
	Id      *uuid.UUID `json:"id,omitempty"`
	Tstamp  *time.Time `json:"tstamp,omitempty"`
}

func (r Relay) Value() (driver.Value, error) {
	b, err := json.Marshal(r)
	return string(b), err
}

func (r Relay) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &r)
}
