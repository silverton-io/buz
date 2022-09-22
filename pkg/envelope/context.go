// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/util"
)

const HTTP_HEADERS_CONTEXT string = "io.silverton/buz/internal/contexts/httpHeaders/v1.0.json"

func BuildContextsFromRequest(c *gin.Context) map[string]interface{} {
	headers := util.HttpHeadersToMap(c)
	context := map[string]interface{}{
		HTTP_HEADERS_CONTEXT: headers,
	}
	return context
}
