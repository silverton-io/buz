package util

import (
	"encoding/json"
)

func Stringify(i interface{}) string {
	payload, _ := json.Marshal(i)
	return string(payload)
}
