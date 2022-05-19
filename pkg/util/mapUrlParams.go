package util

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func MapUrlParams(c *gin.Context) map[string]interface{} {
	// Coerce query params to a map[string]interface{}.
	// Only use the first val of each key.
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
