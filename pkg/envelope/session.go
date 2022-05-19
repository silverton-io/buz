package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Session struct {
	Id     *string                 `json:"id"`
	Idx    *int64                  `json:"idx"`
	Traits *map[string]interface{} `json:"traits"`
}

func (e Session) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e Session) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}
