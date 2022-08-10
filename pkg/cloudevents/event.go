// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package cloudevents

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type CloudEvent struct { // https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/formats/cloudevents.json
	Id              string                 `json:"id"`
	Source          string                 `json:"source"`
	SpecVersion     string                 `json:"specversion"`
	Type            string                 `json:"type"`
	DataContentType string                 `json:"datacontenttype"`
	DataSchema      string                 `json:"dataschema"`
	Subject         *string                `json:"subject"`
	Time            time.Time              `json:"time"`
	Data            map[string]interface{} `json:"data"`
	Datab64         string                 `json:"datab64"`
}

func (e CloudEvent) SchemaName() *string {
	return &e.DataSchema
}

func (e CloudEvent) PayloadAsByte() ([]byte, error) {
	eBytes, err := json.Marshal(e.Data)
	if err != nil {
		return nil, err
	}
	return eBytes, nil
}

func (e CloudEvent) AsByte() ([]byte, error) {
	eBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return eBytes, nil
}

func (e CloudEvent) AsMap() (map[string]interface{}, error) {
	var event map[string]interface{}
	cByte, err := e.AsByte()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(cByte, &event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (e CloudEvent) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e CloudEvent) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}
