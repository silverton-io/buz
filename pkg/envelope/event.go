// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Event interface {
	SchemaName() *string
	PayloadAsByte() ([]byte, error)
	Value() (driver.Value, error)
	Scan(input interface{}) error
}

type Payload map[string]interface{}

func (p Payload) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p Payload) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &p)
}

func (p *Payload) AsByte() ([]byte, error) {
	return json.Marshal(p)
}
