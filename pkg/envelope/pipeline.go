package envelope

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Pipeline struct {
	Source    `json:"source"`
	Collector `json:"collector"`
	Relay     `json:"relay"`
}

func (p Pipeline) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p Pipeline) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &p)
}

type Source struct {
	Ip              string    `json:"ip"`
	GeneratedTstamp time.Time `json:"generatedTstamp"`
	SentTstamp      time.Time `json:"sentTstamp"`
	Name            *string   `json:"name,omitempty"`
	Version         *string   `json:"version,omitempty"`
}

func (s Source) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	return string(b), err
}

func (s Source) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &s)
}

type Collector struct {
	Tstamp time.Time `json:"tstamp"`
}

func (c Collector) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c Collector) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &c)
}

type Relay struct {
	Relayed bool       `json:"relayed"`
	Id      *uuid.UUID `json:"id,omitempty"`
	Tstamp  *time.Time `json:"tstamp,omitempty"`
}

func (r Relay) Value() (driver.Value, error) {
	b, err := json.Marshal(r)
	return string(b), err
}

func (r Relay) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &r)
}
