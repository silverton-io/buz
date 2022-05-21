package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Page struct {
	Page     PageAttrs `json:"page"`
	Referrer PageAttrs `json:"referrer"`
}

func (p Page) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p Page) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &p)
}

type PageAttrs struct {
	Url      string                  `json:"url"`
	Title    *string                 `json:"title"`
	Scheme   string                  `json:"scheme"`
	Host     string                  `json:"host"`
	Port     string                  `json:"port"`
	Path     string                  `json:"path"`
	Query    *map[string]interface{} `json:"query"`
	Fragment *string                 `json:"fragment"`
	Medium   *string                 `json:"medium"`
	Source   *string                 `json:"source"`
	Term     *string                 `json:"term"`
	Content  *string                 `json:"content"`
	Campaign *string                 `json:"campaign"`
}
