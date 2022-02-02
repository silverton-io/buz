package main

import (
	"context"
	"io/ioutil"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func HandleHealthcheck(c *gin.Context) {
	c.JSON(200, AllOk)
}

func HandleRedirect() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		mappedParams := mapParams(c)
		event := buildEventFromMappedParams(c, mappedParams)
		prettyPrint(event)
		// Publish event to Pubsub
		// publishValidEvents([]Event{event})
		redirectUrl, _ := c.GetQuery("u")
		c.Redirect(302, redirectUrl)
	}
	return gin.HandlerFunc(fn)
}

func HandleGet(publishTopic *pubsub.Topic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		mappedParams := mapParams(c)
		event := buildEventFromMappedParams(c, mappedParams)
		// TODO: Validate event
		publishEvent(ctx, publishTopic, event)
		c.JSON(200, AllOk)
	}
	return gin.HandlerFunc(fn)
}

func HandlePost(topic *pubsub.Topic) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		body, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs here
		payloadData := gjson.GetBytes(body, "data")
		var validEvents = []Event{}
		// var invalidEvents = []Event{}
		for _, e := range payloadData.Array() {
			event := buildEventFromMappedParams(c, e.Value().(map[string]interface{}))
			validEvents = append(validEvents, event)
		}
		publishEvents(ctx, topic, validEvents)
		c.JSON(200, AllOk)
	}
	return gin.HandlerFunc(fn)
}
