package request

import "github.com/gin-gonic/gin"

func MapParams(c *gin.Context) map[string]interface{} {
	// Coerce query params to a map[string]interface{}.
	// Only use the first val of each key.
	mappedParams := make(map[string]interface{})
	params := c.Request.URL.Query()
	for k, v := range params {
		mappedParams[k] = v[0]
	}
	return mappedParams
}
