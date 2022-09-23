// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputpixel

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/event"
	"github.com/silverton-io/buz/pkg/util"
)

const (
	B64_ENCODED_PAYLOAD_PARAM string = "hbp"
	ARBITRARY_PIXEL_SCHEMA    string = "io.silverton/buz/pixel/arbitrary/v1.0.json"
)

func buildEvent(c *gin.Context) (event.SelfDescribingPayload, error) {
	params := util.MapUrlParams(c)
	schemaName := util.GetSchemaNameFromRequest(c, ARBITRARY_PIXEL_SCHEMA)
	base64EncodedPayload := params[B64_ENCODED_PAYLOAD_PARAM]
	if base64EncodedPayload != nil {
		p, err := b64.RawStdEncoding.DecodeString(base64EncodedPayload.(string))
		if err != nil {
			return event.SelfDescribingPayload{}, err
		}
		var payload map[string]interface{}
		err = json.Unmarshal(p, &payload)
		if err != nil {
			return event.SelfDescribingPayload{}, err
		}
		sdp := event.SelfDescribingPayload{
			Schema: schemaName,
			Data:   payload,
		}
		return sdp, nil
	}

	sdp := event.SelfDescribingPayload{
		Schema: schemaName,
		Data:   params,
	}
	return sdp, nil
}
