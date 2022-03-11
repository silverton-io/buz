package validator

type PayloadValidationError struct {
	Field       string `json:"field"`
	Description string `json:"description"`
	ErrorType   string `json:"errorType"`
}

type ValidationError struct {
	ErrorType       string                    `json:"errorType"`
	ErrorResolution string                    `json:"errorResolution"`
	Errors          *[]PayloadValidationError `json:"payloadValidationErrors"`
}
