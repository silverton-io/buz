package webhook

import (
	"encoding/json"

	"github.com/silverton-io/honeypot/pkg/protocol"
)

type WebhookEvent struct {
	Payload map[string]interface{} `json:"payload"`
}

func (e WebhookEvent) Schema() *string {
	schema := "" // FIXME! Should incoming webhooks have an associated schema? IDK yet.
	return &schema
}

func (e WebhookEvent) Protocol() string {
	return protocol.WEBHOOK
}

func (e WebhookEvent) PayloadAsByte() ([]byte, error) {
	payloadBytes, err := json.Marshal(e.Payload)
	if err != nil {
		return nil, err
	}
	return payloadBytes, nil
}

func (e WebhookEvent) AsByte() ([]byte, error) {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return eventBytes, nil
}

func (e WebhookEvent) AsMap() (map[string]interface{}, error) {
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
