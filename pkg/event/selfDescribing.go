package event

type SelfDescribingEnvelope struct {
	Contexts []SelfDescribingContext `json:"contexts"`
	Payload  SelfDescribingPayload   `json:"payload"`
}

type SelfDescribingPayload struct {
	Schema string                 `json:"schema"`
	Data   map[string]interface{} `json:"data"`
}

type SelfDescribingContext SelfDescribingPayload
