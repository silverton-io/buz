// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package snowplow

import (
	"testing"

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

// FIXME
func TestSnowplowEvent(t *testing.T) {}

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
