package snowplow

import (
	"context"
	"io/ioutil"

	"github.com/gin-gonic/gin"

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
		isValid, validationError, schema := ValidateEvent(event, cache)
		setEventMetadataFields(&event, schema)
		if isValid {
			forwarder.PublishValidEvent(ctx, event)
		} else {
			invalidEvent := InvalidEvent{
				ValidationError: &validationError,
				Event:           &event,
			}
			forwarder.PublishInvalidEvent(ctx, invalidEvent)
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
		var invalidEvents []interface{}
		for _, e := range payloadData.Array() {
			event := BuildEventFromMappedParams(c, e.Value().(map[string]interface{}))
			isValid, validationError, schema := ValidateEvent(event, cache)
			setEventMetadataFields(&event, schema)
			if isValid {
				validEvents = append(validEvents, event)
			} else {
				invalidEvent := InvalidEvent{
					ValidationError: &validationError,
					Event:           &event,
				}
				invalidEvents = append(invalidEvents, invalidEvent)
			}
		}
		forwarder.PublishValidEvents(ctx, validEvents)
		forwarder.PublishInvalidEvents(ctx, invalidEvents)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
