// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type Web struct {
	Page     PageAttrs `json:"page"`
	Referrer PageAttrs `json:"referrer"`
}

func (p Web) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p Web) Scan(input interface{}) error {
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
