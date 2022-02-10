package handler

import (
	"context"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/http"
	"github.com/silverton-io/gosnowplow/pkg/response"
	"github.com/silverton-io/gosnowplow/pkg/snowplow"
	"github.com/tidwall/gjson"
)

func Healthcheck(c *gin.Context) {
	c.JSON(200, response.Ok)
}

func SnowplowRedirect(forwarder *forwarder.PubsubForwarder) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		mappedParams := http.MapParams(c)
		event := snowplow.BuildEventFromMappedParams(c, mappedParams)
		forwarder.PublishValidEvent(ctx, event)
		redirectUrl, _ := c.GetQuery("u")
		c.Redirect(302, redirectUrl)
	}
	return gin.HandlerFunc(fn)
}

func SnowplowGet(forwarder *forwarder.PubsubForwarder) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		mappedParams := http.MapParams(c)
		event := snowplow.BuildEventFromMappedParams(c, mappedParams)
		forwarder.PublishValidEvent(ctx, event)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}

func SnowplowPost(forwarder *forwarder.PubsubForwarder) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		body, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs here
		payloadData := gjson.GetBytes(body, "data")
		var validEvents []interface{}
		for _, e := range payloadData.Array() {
			event := snowplow.BuildEventFromMappedParams(c, e.Value().(map[string]interface{}))
			validEvents = append(validEvents, event)
		}
		forwarder.PublishValidEvents(ctx, validEvents)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
