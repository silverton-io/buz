package util

import (
	"encoding/json"
	"fmt"
)

func PrintByteSliceAsJson(b []byte) {
	payload := make(map[string]interface{})
	json.Unmarshal(b, &payload)
	PrettyPrint(payload)
}

func PrettyPrint(i interface{}) {
	payload, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(payload))
}

func Stringify(i interface{}) string {
	payload, _ := json.Marshal(i)
	return string(payload)
}
