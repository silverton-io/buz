// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"database/sql/driver"
	"encoding/json"
)

type User struct {
	Id          *string                 `json:"id,omitempty"`
	AnonymousId *string                 `json:"anonymousId,omitempty"`
	Fingerprint *string                 `json:"fingerprint,omitempty"`
	Traits      *map[string]interface{} `json:"traits,omitempty"`
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
	Name *string `json:"name,omitempty"`
}

type Group struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}
