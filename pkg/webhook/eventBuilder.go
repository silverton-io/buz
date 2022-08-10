// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/tidwall/gjson"
)

const ARBITRARY_WEBHOOK_SCHEMA = "io.silverton/honeypot/hook/arbitrary/v1.0.json"

func BuildEvent(c *gin.Context, payload gjson.Result) (event.SelfDescribingPayload, error) {
	schemaName := util.GetSchemaNameFromRequest(c, ARBITRARY_WEBHOOK_SCHEMA)
	sdp := event.SelfDescribingPayload{
		Schema: schemaName,
		Data:   payload.Value().(map[string]interface{}),
	}
	return sdp, nil
}
