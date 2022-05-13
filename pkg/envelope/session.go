package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Session struct {
	Sid  *string `json:"sid"`
	Sidx *int64  `json:"sidx"`
}

func (e Session) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e Session) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}
