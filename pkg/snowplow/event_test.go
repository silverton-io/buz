package snowplow

import (
	"encoding/json"
	"testing"

	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, PAGE_PING, "page_ping")
	assert.Equal(t, PAGE_VIEW, "page_view")
	assert.Equal(t, STRUCT_EVENT, "struct_event")
	assert.Equal(t, SELF_DESCRIBING_EVENT, "self_describing")
	assert.Equal(t, TRANSACTION, "transaction")
	assert.Equal(t, TRANSACTION_ITEM, "transaction_item")
	assert.Equal(t, AD_IMPRESSION, "ad_impression")
	assert.Equal(t, UNKNOWN_EVENT, "unknown_event")
	assert.Equal(t, UNKNOWN_SCHEMA, "unknown_schema")
}

func TestSnowplowEvent(t *testing.T) {
	schema := "com.something.somewhere/did/v1.0.json"
	gooddata := map[string]interface{}{
		"something": "somewhere",
		"when":      "then",
		"count":     10,
	}
	sdPayload := event.SelfDescribingPayload{
		Schema: schema,
		Data:   gooddata,
	}
	e := SnowplowEvent{
		SelfDescribingEvent: &sdPayload,
	}
	expectedPayloadByte, _ := json.Marshal(e.SelfDescribingEvent.Data)
	actualPayloadByte, _ := e.PayloadAsByte()
	expectedByte, _ := json.Marshal(e)
	actualByte, _ := e.AsByte()
	var expectedMap map[string]interface{}
	actualMap, _ := e.AsMap()
	json.Unmarshal(actualByte, &expectedMap)
	assert.Equal(t, schema, *e.Schema())
	assert.Equal(t, protocol.SNOWPLOW, e.Protocol())
	assert.Equal(t, expectedPayloadByte, actualPayloadByte)
	assert.Equal(t, expectedByte, actualByte)
	assert.Equal(t, expectedMap, actualMap)
}

func TestGetEventType(t *testing.T) {
	assert.Equal(t, PAGE_PING, getEventType("pp"))
	assert.Equal(t, PAGE_VIEW, getEventType("pv"))
	assert.Equal(t, STRUCT_EVENT, getEventType("se"))
	assert.Equal(t, SELF_DESCRIBING_EVENT, getEventType("ue"))
	assert.Equal(t, TRANSACTION, getEventType("tr"))
	assert.Equal(t, TRANSACTION_ITEM, getEventType("ti"))
	assert.Equal(t, AD_IMPRESSION, getEventType("ad"))
	assert.Equal(t, UNKNOWN_EVENT, getEventType("yikes"))
}
