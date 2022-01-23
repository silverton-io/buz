package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type SelfDescribingPayload struct {
	Schema string                 `json:"schema"`
	Data   map[string]interface{} `json:"data"`
}

type Base64EncodedSelfDescribingPayload struct {
	Schema string                 `json:"schema"`
	Data   map[string]interface{} `json:"data"`
}

func (f *Base64EncodedSelfDescribingPayload) UnmarshalJSON(bytes []byte) error {
	var encodedPayload string
	err := json.Unmarshal(bytes, &encodedPayload)
	decodedPayload, err := b64.RawStdEncoding.DecodeString(encodedPayload)
	if err != nil {
		fmt.Printf("error decoding b64 encoded self describing payload %s\n", err)
	}
	schema := gjson.GetBytes(decodedPayload, "data.schema").String()
	data := gjson.GetBytes(decodedPayload, "data.data").Value().(map[string]interface{})
	fmt.Println("SCHEMA: ", schema)
	fmt.Println("DATA: ", data)
	*&f.Schema = schema
	*&f.Data = data
	return nil
}

type MillisecondTimestampField struct {
	time.Time
}

func (t *MillisecondTimestampField) UnmarshalJSON(bytes []byte) error {
	var msString string
	err := json.Unmarshal(bytes, &msString)
	msInt, err := strconv.ParseInt(msString, 10, 64)
	if err != nil {
		fmt.Printf("error decoding timestamp: %s\n", err)
		return err
	}
	*&t.Time = time.Unix(0, msInt*int64(time.Millisecond))
	return nil
}

// func stringToWidth(val string) int {
// 	split := strings.Split(val, "x")
// 	return stringToInt(split[0])
// }

// func stringToHeight(val string) int {
// 	split := strings.Split(val, "x")
// 	return stringToInt(split[1])
// }

func getEventType(param string) string {
	switch param {
	case "pp":
		return "page_ping"
	case "pv":
		return "page_view"
	case "se":
		return "struct_event"
	case "ue":
		return "self_describing"
	case "tr":
		return "transaction"
	case "ti":
		return "transaction_item"
	case "ad":
		return "ad_impression"
	}
	return "unknown"
}
