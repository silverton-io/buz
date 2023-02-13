// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/event"
	"github.com/silverton-io/buz/pkg/util"
	"github.com/tidwall/gjson"
)

const ARBITRARY_WEBHOOK_SCHEMA = "io.silverton/buz/hook/arbitrary/v1.0.json"

func buildEvent(c *gin.Context, payload gjson.Result) (event.SelfDescribingPayload, error) {
	schemaName := util.GetSchemaNameFromRequest(c, ARBITRARY_WEBHOOK_SCHEMA)
	sdp := event.SelfDescribingPayload{
		Schema: schemaName,
		Data:   payload.Value().(map[string]interface{}),
	}
	return sdp, nil
}
