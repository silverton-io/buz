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
	Url      string                  `json:"url,omitempty"`
	Title    *string                 `json:"title,omitempty"`
	Scheme   string                  `json:"scheme,omitempty"`
	Host     string                  `json:"host,omitempty"`
	Port     string                  `json:"port,omitempty"`
	Path     string                  `json:"path,omitempty"`
	Query    *map[string]interface{} `json:"query,omitempty"`
	Fragment *string                 `json:"fragment,omitempty"`
	Medium   *string                 `json:"medium,omitempty"`
	Source   *string                 `json:"source,omitempty"`
	Term     *string                 `json:"term,omitempty"`
	Content  *string                 `json:"content,omitempty"`
	Campaign *string                 `json:"campaign,omitempty"`
}
