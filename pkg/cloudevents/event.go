package cloudevents

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/protocol"
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
	DataBase64      string                 `json:"dataBase64"`
}

func (e CloudEvent) Schema() *string {
	return &e.DataSchema
}

func (e CloudEvent) Protocol() string {
	return protocol.CLOUDEVENTS
}

func (e CloudEvent) AsByte() ([]byte, error) {
	eBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not marshal cloudevent")
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
