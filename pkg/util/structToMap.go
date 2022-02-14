package util

import "encoding/json"

func StructToMap(v interface{}) map[string]interface{} {
	var m map[string]interface{}
	i, _ := json.Marshal(v)
	json.Unmarshal(i, &m)
	return m
}
