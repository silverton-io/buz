package event

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

type SelfDescribingEvent struct {
	Contexts []SelfDescribingContext `json:"contexts"`
	Payload  SelfDescribingPayload   `json:"payload"`
}

type SelfDescribingPayload struct {
	Schema string                 `json:"schema"`
	Data   map[string]interface{} `json:"data"`
}

func (f *SelfDescribingPayload) UnmarshalJSON(bytes []byte) error {
	var encodedPayload string

	err := json.Unmarshal(bytes, &encodedPayload)
	decodedPayload, err := b64.RawStdEncoding.DecodeString(encodedPayload)
	schema := gjson.GetBytes(decodedPayload, "data.schema").String()
	data := gjson.GetBytes(decodedPayload, "data.data").Value().(map[string]interface{})
	fmt.Print("\n\n")
	fmt.Println(schema)
	fmt.Println(data)
	fmt.Print("\n\n")
	f.Schema = schema
	f.Data = data
	if err != nil {
		log.Error().Err(err).Msg("error decoding b64 encoded self describing payload")
		return err
	}
	// // First assume it's json
	// schema := gjson.GetBytes(bytes, "data.schema").String()
	// data := gjson.GetBytes(bytes, "data.data").Value()
	// if data != nil {
	// 	f.Schema = schema
	// 	f.Data = data.(map[string]interface{})
	// } else {
	return nil
}

type SelfDescribingContext SelfDescribingPayload

type SelfDescribingContexts []SelfDescribingContext

func (c *SelfDescribingContexts) UnmarshalJSON(bytes []byte) error {
	var encodedPayload string
	var contexts []SelfDescribingContext
	err := json.Unmarshal(bytes, &encodedPayload)
	decodedPayload, err := b64.RawStdEncoding.DecodeString(encodedPayload)
	if err != nil {
		log.Error().Err(err).Msg("error decoding b64 encoded contexts")
		log.Fatal().Err(err)
	}
	contextPayload := gjson.Parse(string(decodedPayload))
	for _, pl := range contextPayload.Get("data").Array() {
		context := SelfDescribingContext{
			Schema: pl.Get("schema").String(),
			Data:   pl.Get("data").Value().(map[string]interface{}),
		}
		contexts = append(contexts, context)
	}
	*c = contexts
	return nil
}
