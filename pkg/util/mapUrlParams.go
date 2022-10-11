// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

// Coerce query params to a map[string]interface{}.
// Only use the first val of each key.
func MapUrlParams(c *gin.Context) map[string]interface{} {
	mappedParams := make(map[string]interface{})
	params := c.Request.URL.Query()
	for k, v := range params {
		mappedParams[k] = v[0]
	}
	return mappedParams
}

func QueryToMap(v url.Values) map[string]interface{} {
	mappedParams := make(map[string]interface{})
	for k, v := range v {
		mappedParams[k] = v[0]
	}
	return mappedParams
}
