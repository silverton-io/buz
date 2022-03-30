package snowplow

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		height: 100,
		width:  200,
	}
	actual, _ := getDimensions(dimString)
	assert.Equal(t, expected, actual)
}

func TestGetContexts(t *testing.T) {

}

func TestGetSdPayload(t *testing.T) {

}

func TestGetQueryParam(t *testing.T) {

}

func TestGetPageFieldsFromUrl(t *testing.T) {

}

func TestSetTsFields(t *testing.T) {

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
