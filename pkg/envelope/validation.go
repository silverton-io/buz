package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Validation struct {
	IsValid bool             `json:"isValid"`
	Error   *ValidationError `json:"error"`
}

func (e Validation) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e Validation) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}

type PayloadValidationError struct {
	Field       string `json:"field"`
	Description string `json:"description"`
	ErrorType   string `json:"errorType"`
}

type ValidationError struct {
	ErrorType       *string                  `json:"errorType"`
	ErrorResolution *string                  `json:"errorResolution"`
	Errors          []PayloadValidationError `json:"payloadValidationErrors"`
}

func (e *ValidationError) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e *ValidationError) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), e)
}
