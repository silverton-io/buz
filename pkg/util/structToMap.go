// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"encoding/json"
)

func StructToMap(v interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}
	i, _ := json.Marshal(v)
	if err := json.Unmarshal(i, &m); err != nil {
		return map[string]interface{}{}, err
	}
	return m, nil
}
