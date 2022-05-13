package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type User struct {
	Uid         *string `json:"uid"`
	Fingerprint *string `json:"fingerprint"`
	Groups      []Group `json:"groups"`
}

func (e User) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e User) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}

type Group struct {
	Name *string `json:"name"`
}
