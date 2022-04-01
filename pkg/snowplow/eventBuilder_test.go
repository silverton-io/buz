package snowplow

import (
	b64 "encoding/base64"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func buildMockMap() map[string]interface{} {
	mm := map[string]interface{}{
		"s": "somestring",
		"f": "23.99",
		"i": "10",
		"b": "true",
		"t": "1648667060951",
	}
	return mm
}

func TestGetStringParam(t *testing.T) {
	mm := buildMockMap()
	expected := "somestring"
	actual := getStringParam(mm, "s")
	assert.Equal(t, expected, *actual)
}

func TestGetInt64Param(t *testing.T) {
	mm := buildMockMap()
	var expected int64 = 10
	actual := getInt64Param(mm, "i")
	assert.Equal(t, expected, *actual)
}

func TestGetFloat64Param(t *testing.T) {
	mm := buildMockMap()
	var expected float64 = 23.99
	actual := getFloat64Param(mm, "f")
	assert.Equal(t, expected, *actual)
}

func TestGetTimeParam(t *testing.T) {
	mm := buildMockMap()
	iVal, _ := strconv.ParseInt("1648667060951", 10, 64)
	expected := time.UnixMilli(iVal)
	actual := getTimeParam(mm, "t")
	assert.Equal(t, expected, *actual)
}

func TestGetBoolParam(t *testing.T) {
	mm := buildMockMap()
	var expected bool = true
	actual := getBoolParam(mm, "b")
	assert.Equal(t, expected, *actual)
}

func TestGetDimensions(t *testing.T) {
	dimString := "100x200"
	expected := Dimension{
		width:  100,
		height: 200,
	}
	actual, _ := getDimensions(dimString)
	assert.Equal(t, expected, actual)
}

func TestGetContexts(t *testing.T) {
	b64contexts := "eyJzY2hlbWEiOiJpZ2x1OmNvbS5zbm93cGxvd2FuYWx5dGljcy5zbm93cGxvdy9jb250ZXh0cy9qc29uc2NoZW1hLzEtMC0wIiwiZGF0YSI6W3sic2NoZW1hIjoiaWdsdTpjb20uc25vd3Bsb3dhbmFseXRpY3Muc25vd3Bsb3cvd2ViX3BhZ2UvanNvbnNjaGVtYS8xLTAtMCIsImRhdGEiOnsiaWQiOiI0ZTRjM2UzMS05Y2FkLTQ1YjgtYTMzOC1kMzNiN2E4ODQwMzQifX0seyJzY2hlbWEiOiJpZ2x1Om9yZy53My9QZXJmb3JtYW5jZVRpbWluZy9qc29uc2NoZW1hLzEtMC0wIiwiZGF0YSI6eyJuYXZpZ2F0aW9uU3RhcnQiOjE2NDg2NzEwOTQ1MTksInJlZGlyZWN0U3RhcnQiOjAsInJlZGlyZWN0RW5kIjowLCJmZXRjaFN0YXJ0IjoxNjQ4NjcxMDk3MDk3LCJkb21haW5Mb29rdXBTdGFydCI6MTY0ODY3MTA5NzEwMiwiZG9tYWluTG9va3VwRW5kIjoxNjQ4NjcxMDk3MTAyLCJjb25uZWN0U3RhcnQiOjE2NDg2NzEwOTcxMDIsInNlY3VyZUNvbm5lY3Rpb25TdGFydCI6MCwiY29ubmVjdEVuZCI6MTY0ODY3MTA5NzEwMywicmVxdWVzdFN0YXJ0IjoxNjQ4NjcxMDk3MTAzLCJyZXNwb25zZVN0YXJ0IjoxNjQ4NjcxMDk3MTA3LCJyZXNwb25zZUVuZCI6MTY0ODY3MTA5NzEwNywidW5sb2FkRXZlbnRTdGFydCI6MTY0ODY3MTA5NzExMCwidW5sb2FkRXZlbnRFbmQiOjE2NDg2NzEwOTcxMTAsImRvbUxvYWRpbmciOjE2NDg2NzEwOTQ1MjAsImRvbUludGVyYWN0aXZlIjoxNjQ4NjcxMDk0NTMyLCJkb21Db250ZW50TG9hZGVkRXZlbnRTdGFydCI6MTY0ODY3MTA5NDU3MywiZG9tQ29udGVudExvYWRlZEV2ZW50RW5kIjoxNjQ4NjcxMDk0NTc0LCJkb21Db21wbGV0ZSI6MTY0ODY3MTA5OTg4OSwibG9hZEV2ZW50U3RhcnQiOjE2NDg2NzEwOTk4ODksImxvYWRFdmVudEVuZCI6MTY0ODY3MTA5OTg4OX19XX0"
	var expectedContexts []event.SelfDescribingContext
	pl, _ := b64.RawStdEncoding.DecodeString(b64contexts)
	contextPayload := gjson.ParseBytes(pl)
	for _, pl := range contextPayload.Get("data").Array() {
		c := event.SelfDescribingContext{
			Schema: pl.Get("schema").String(),
			Data:   pl.Get("data").Value().(map[string]interface{}),
		}
		expectedContexts = append(expectedContexts, c)
	}
	actualContext := getContexts(&b64contexts)
	assert.Equal(t, expectedContexts, *actualContext)
}

func TestGetSdPayload(t *testing.T) {
	b64payload := "eyJzY2hlbWEiOiJpZ2x1OmNvbS5zbm93cGxvd2FuYWx5dGljcy5zbm93cGxvdy91bnN0cnVjdF9ldmVudC9qc29uc2NoZW1hLzEtMC0wIiwiZGF0YSI6eyJzY2hlbWEiOiJpZ2x1OmNvbS5zaWx2ZXJ0b24uaW8vaG9uZXlwb3QvZXhhbXBsZS92aWV3ZWRfcHJvZHVjdC9qc29uc2NoZW1hLzEtMC0wIiwiZGF0YSI6eyJwcm9kdWN0SWQiOiJBU08wMTA0MyIsImNhdGVnb3J5IjoiRHJlc3NlcyIsImJyYW5kIjoiQUNNRSIsInJldHVybmluZyI6dHJ1ZSwicHJpY2UiOjQ5Ljk1LCJzaXplcyI6WyJ4cyIsInMiLCJsIiwieGwiLCJ4eGwiXSwiYXZhaWxhYmxlU2luY2UiOiIyMDEzLTA0LTA3VDA0OjAwOjAwLjAwMFoifX19"
	pl, _ := b64.RawStdEncoding.DecodeString(b64payload)
	payload := gjson.ParseBytes(pl)
	expectedPayload := event.SelfDescribingPayload{
		Schema: payload.Get("data.schema").String(),
		Data:   payload.Get("data.data").Value().(map[string]interface{}),
	}
	actualPayload := getSdPayload(&b64payload)
	assert.Equal(t, expectedPayload, *actualPayload)
}

func TestGetQueryParam(t *testing.T) {
	u, _ := url.Parse("http://somewhere.net?q=100")
	v1 := "100"
	var v2 *string
	p1 := getQueryParam(*u, "q")
	p2 := getQueryParam(*u, "s")
	assert.Equal(t, v1, *p1)
	assert.Equal(t, v2, p2)
}

func TestGetPageFieldsFromUrl(t *testing.T) {
}

func TestSetMetadataFields(t *testing.T) {
}

func TestSetUserFields(t *testing.T) {

}

func TestSetBrowserFeatures(t *testing.T) {

}

func TestSetDimensionFields(t *testing.T) {

}

func TestSetPageFields(t *testing.T) {

}

func TestSetReferrerFields(t *testing.T) {}

func TestAnonymizeFields(t *testing.T) {

}

func TestSetPagePingFields(t *testing.T) {

}

func TestSetStructFields(t *testing.T) {

}

func TestSetTransactionFields(t *testing.T) {

}

func TestSetTransactionItemFields(t *testing.T) {

}

func TestSetContexts(t *testing.T) {

}

func TestSetSelfDescribingFields(t *testing.T) {

}

func TestSetEventMetadataFields(t *testing.T) {

}

func TestBuildEventFromMappedParams(t *testing.T) {

}
