// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Validation struct {
	IsValid *bool            `json:"isValid"`
	Error   *ValidationError `json:"error,omitempty"`
}

func (e Validation) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e Validation) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}

type PayloadValidationError struct {
	Field       string `json:"field,omitempty"`
	Description string `json:"description,omitempty"`
	ErrorType   string `json:"errorType,omitempty"`
}

type ValidationError struct {
	ErrorType       *string                  `json:"errorType,omitempty"`
	ErrorResolution *string                  `json:"errorResolution,omitempty"`
	Errors          []PayloadValidationError `json:"payloadValidationErrors,omitempty"`
}

func (e *ValidationError) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e *ValidationError) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), e)
}
