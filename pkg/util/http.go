// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import "github.com/gin-gonic/gin"

func HttpHeadersToMap(c *gin.Context) map[string]interface{} {
	headers := make(map[string]interface{})
	for k, v := range c.Request.Header {
		if len(v) == 1 {
			headers[k] = v[0]
		} else {
			headers[k] = v
		}
	}
	return headers
}
