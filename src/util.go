package main

import (
	"encoding/json"
	"fmt"
)

func prettyPrint(i interface{}) {
	payload, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(payload))
}
