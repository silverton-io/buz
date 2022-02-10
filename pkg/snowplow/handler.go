package snowplow

import (
	"context"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/http"
	"github.com/silverton-io/gosnowplow/pkg/response"
	"github.com/tidwall/gjson"
)

func RedirectHandler(forwarder *forwarder.PubsubForwarder, cache *cache.SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		mappedParams := http.MapParams(c)
		event := BuildEventFromMappedParams(c, mappedParams)
		forwarder.PublishValidEvent(ctx, event)
		redirectUrl, _ := c.GetQuery("u")
		c.Redirect(302, redirectUrl)
	}
	return gin.HandlerFunc(fn)
}

func GetHandler(forwarder *forwarder.PubsubForwarder, cache *cache.SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		mappedParams := http.MapParams(c)
		event := BuildEventFromMappedParams(c, mappedParams)
		isValid, err := ValidateEvent(event, cache)
		if isValid {
			forwarder.PublishValidEvent(ctx, event)
		} else {
			forwarder.PublishInvalidEvent(ctx, event)
		}
		if err != nil {
			log.Error().Stack().Msg("error when validating event")
		}
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}

func PostHandler(forwarder *forwarder.PubsubForwarder, cache *cache.SchemaCache) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		body, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs here
		payloadData := gjson.GetBytes(body, "data")
		var validEvents []interface{}
		for _, e := range payloadData.Array() {
			event := BuildEventFromMappedParams(c, e.Value().(map[string]interface{}))
			validEvents = append(validEvents, event)
		}
		forwarder.PublishValidEvents(ctx, validEvents)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
