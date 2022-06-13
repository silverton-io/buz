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
