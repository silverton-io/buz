// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"encoding/json"
	"fmt"
)

func Pprint(i interface{}) string {
	payload, _ := json.MarshalIndent(i, "", "\t")
	stringified := string(payload)
	fmt.Println(stringified)
	return stringified
}

func Stringify(i interface{}) string {
	payload, _ := json.Marshal(i)
	return string(payload)
}
