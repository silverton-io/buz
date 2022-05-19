package pixel

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/silverton-io/honeypot/pkg/protocol"
)

type PixelEvent struct {
	Id      string                 `json:"id"`
	Payload map[string]interface{} `json:"payload"`
}

func (e PixelEvent) SchemaName() *string {
	schema := "" // FIXME! Should incoming pixels associated schema? IDK yet.
	return &schema
}

func (e PixelEvent) Protocol() string {
	return protocol.PIXEL
}

func (e PixelEvent) PayloadAsByte() ([]byte, error) {
	payloadBytes, err := json.Marshal(e.Payload)
	if err != nil {
		return nil, err
	}
	return payloadBytes, nil
}

func (e PixelEvent) AsByte() ([]byte, error) {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return eventBytes, nil
}

func (e PixelEvent) AsMap() (map[string]interface{}, error) {
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

func (e PixelEvent) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e PixelEvent) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}
