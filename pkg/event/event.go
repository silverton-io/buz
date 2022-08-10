// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package event

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

func (p *Payload) AsByte() ([]byte, error) {
	return json.Marshal(p)
}
