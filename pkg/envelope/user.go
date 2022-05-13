package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type User struct {
	Id          *string `json:"id,omitempty"`
	AnonymousId *string `json:"anonymousId,omitempty"`
	Fingerprint *string `json:"fingerprint,omitempty"`
}

func (e User) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e User) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}

type Team struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name"`
}

type Group struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name"`
}
