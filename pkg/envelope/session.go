package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Session struct {
	Id     *string                 `json:"id,omitempty"`
	Idx    *int64                  `json:"idx,omitempty"`
	Traits *map[string]interface{} `json:"traits,omitempty"`
}

func (e Session) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e Session) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}
