package util

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(i interface{}) {
	payload, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(payload))
}
