package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type User struct {
	Id          *string                 `json:"id"`
	AnonymousId *string                 `json:"anonymousId"`
	Fingerprint *string                 `json:"fingerprint"`
	Traits      *map[string]interface{} `json:"traits"`
}

func (e User) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e User) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}

type Team struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type Group struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}
