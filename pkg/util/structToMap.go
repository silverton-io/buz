// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import "encoding/json"

func StructToMap(v interface{}) map[string]interface{} {
	var m map[string]interface{}
	i, _ := json.Marshal(v)
	// FIXME! Don't love it. Don't love it at all.
	if err := json.Unmarshal(i, &m); err != nil {
		panic(err)
	}
	return m
}
