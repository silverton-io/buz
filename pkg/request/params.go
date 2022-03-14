package request

import "github.com/gin-gonic/gin"

func MapParams(c *gin.Context) map[string]string {
	// Coerce query params to a map[string]string.
	// Only use the first val of each key.
	mappedParams := make(map[string]string)
	params := c.Request.URL.Query()
	for k, v := range params {
		mappedParams[k] = v[0]
	}
	return mappedParams
}
