package handler

import (
	"context"
	"io/ioutil"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"

	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/http"
	"github.com/silverton-io/gosnowplow/pkg/response"
	"github.com/silverton-io/gosnowplow/pkg/snowplow"
	"github.com/silverton-io/gosnowplow/pkg/util"
	"github.com/tidwall/gjson"
)

func Healthcheck(c *gin.Context) {
	c.JSON(200, response.Ok)
}

func SnowplowRedirect(validTopic *pubsub.Topic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		mappedParams := http.MapParams(c)
		event := snowplow.BuildEventFromMappedParams(c, mappedParams)
		forwarder.PublishEvent(ctx, validTopic, event)
		redirectUrl, _ := c.GetQuery("u")
		c.Redirect(302, redirectUrl)
	}
	return gin.HandlerFunc(fn)
}

func SnowplowGet(publishTopic *pubsub.Topic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		mappedParams := http.MapParams(c)
		event := snowplow.BuildEventFromMappedParams(c, mappedParams)
		util.PrettyPrint(event)
		forwarder.PublishEvent(ctx, publishTopic, event)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}

func SnowplowPost(topic *pubsub.Topic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		body, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs here
		payloadData := gjson.GetBytes(body, "data")
		var validEvents = []snowplow.Event{}
		for _, e := range payloadData.Array() {
			event := snowplow.BuildEventFromMappedParams(c, e.Value().(map[string]interface{}))
			validEvents = append(validEvents, event)
		}
		forwarder.PublishEvents(ctx, topic, validEvents)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
