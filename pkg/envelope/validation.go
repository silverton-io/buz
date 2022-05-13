package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type ValidationMeta struct {
	IsValid         bool             `json:"isValid"`
	ValidationError *ValidationError `json:"validationErrors"`
}

func (e ValidationMeta) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e ValidationMeta) Scan(input interface{}) error {
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
