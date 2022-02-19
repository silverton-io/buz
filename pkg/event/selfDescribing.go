package event

type SelfDescribingEnvelope struct {
	Contexts []SelfDescribingContext `json:"contexts"`
	Event    SelfDescribingPayload   `json:"event"`
}

type SelfDescribingPayload struct {
	Schema string                 `json:"schema"`
	Data   map[string]interface{} `json:"data"`
}

type SelfDescribingContext SelfDescribingPayload
