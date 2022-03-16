package generic

import (
	"encoding/json"

	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
)

type GenericEvent event.SelfDescribingEvent

func (e GenericEvent) Schema() *string {
	return &e.Payload.Schema
}

func (e GenericEvent) Protocol() string {
	return protocol.GENERIC
}

func (e GenericEvent) PayloadAsByte() ([]byte, error) {
	payloadBytes, err := json.Marshal(e.Payload.Data)
	if err != nil {
		return nil, err
	}
	return payloadBytes, nil
}

func (e GenericEvent) AsByte() ([]byte, error) {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return eventBytes, nil
}

func (e GenericEvent) AsMap() (map[string]interface{}, error) {
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
