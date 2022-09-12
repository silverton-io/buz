// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package event

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

func stripColonSeparatedPrefix(schema string) string {
	colonIdx := strings.Index(schema, ":")
	if colonIdx == -1 {
		return schema
	} else {
		return schema[colonIdx+1:]
	}
}

type SelfDescribingEvent struct {
	Contexts map[string]interface{} `json:"contexts"`
	Payload  SelfDescribingPayload  `json:"payload"`
}

type SelfDescribingPayload struct {
	Schema string                 `json:"schema"`
	Data   map[string]interface{} `json:"data"`
}

func (e SelfDescribingPayload) SchemaName() *string {
	name := stripColonSeparatedPrefix(e.Schema)
	return &name
}

func (e SelfDescribingPayload) PayloadAsByte() ([]byte, error) {
	payloadBytes, err := json.Marshal(e.Data)
	if err != nil {
		return nil, err
	}
	return payloadBytes, nil
}

func (e SelfDescribingPayload) AsByte() ([]byte, error) {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return eventBytes, nil
}

func (e SelfDescribingPayload) AsMap() (map[string]interface{}, error) {
	var m map[string]interface{}
	b, err := e.AsByte()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (e SelfDescribingPayload) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e SelfDescribingPayload) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}

type SelfDescribingContext SelfDescribingPayload
