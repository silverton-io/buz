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
	Name    string `json:"name"`
	Version string `json:"version"`
	Ip      string `json:"ip"`
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
	IsRelayed bool       `json:"isRelayed"`
	RelayId   *uuid.UUID `json:"relayId"`
	Tstamp    *time.Time `json:"tstamp"`
}

func (r Relay) Value() (driver.Value, error) {
	b, err := json.Marshal(r)
	return string(b), err
}

func (r Relay) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &r)
}
