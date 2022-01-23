package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type SelfDescribingPayload struct {
	Schema string                 `json:"schema"`
	Data   map[string]interface{} `json:"data"`
}

type Context SelfDescribingPayload

type Base64EncodedContexts []Context

type Base64EncodedSelfDescribingPayload SelfDescribingPayload

func (f *Base64EncodedSelfDescribingPayload) UnmarshalJSON(bytes []byte) error {
	var encodedPayload string
	err := json.Unmarshal(bytes, &encodedPayload)
	decodedPayload, err := b64.RawStdEncoding.DecodeString(encodedPayload)
	if err != nil {
		fmt.Printf("error decoding b64 encoded self describing payload %s\n", err)
	}
	schema := gjson.GetBytes(decodedPayload, "data.schema").String()
	data := gjson.GetBytes(decodedPayload, "data.data").Value().(map[string]interface{})
	*&f.Schema = schema
	*&f.Data = data
	return nil
}

type FlexibleBoolField bool

func (b *FlexibleBoolField) UnmarshalJSON(bytes []byte) error {
	var payload string
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		fmt.Printf("error decoding FlexibleBoolField %s\n", err)
	}
	val, err := strconv.ParseBool(payload)
	*b = FlexibleBoolField(val)
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

type WidthField int

func (wf *WidthField) UnmarshalJSON(bytes []byte) error {
	var payload string
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		fmt.Printf("error unmarshalling WidthField %s\n", err)
	}
	w := strings.Split(payload, "x")[0]
	width, err := strconv.Atoi(w)
	if err != nil {
		fmt.Printf("error converting width string to int")
	}
	*wf = WidthField(width)
	return nil
}

type HeightField int

func (hf *HeightField) UnmarshalJSON(bytes []byte) error {
	var payload string
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		fmt.Printf("error unmarshalling HeightField %s\n", err)
	}
	h := strings.Split(payload, "x")[1]
	height, err := strconv.Atoi(h)
	if err != nil {
		fmt.Printf("error converting height string to int")
	}
	*hf = HeightField(height)
	return nil
}

type EventTypeField string

func (f *EventTypeField) UnmarshalJSON(bytes []byte) error {
	var payload string
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		fmt.Printf("error unmarshaling EventTypeField %s\n", err)
	}
	eventType := getEventType(payload)
	*f = EventTypeField(eventType)
	return nil
}

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

type Dimension struct {
	height int
	width  int
}
