package ce

import (
	"io/ioutil"
	"time"

	c "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	e "github.com/silverton-io/gosnowplow/pkg/event"
	f "github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/response"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"github.com/tidwall/gjson"
)

func bifurcateEvents(events []event.Event, cache *cache.SchemaCache) (validEvents []interface{}, invalidEvents []interface{}) {
	var vEvents []interface{}
	var invEvents []interface{}
	for _, event := range events {
		isValid, validationError := validateEvent(event, cache)
		if isValid {
			vEvents = append(vEvents, event)
		} else {
			invalidEvent := e.InvalidEvent{
				ValidationError: &validationError,
				Event:           event,
			}
			invEvents = append(invEvents, invalidEvent)
		}
	}
	return vEvents, invEvents
}

func buildCloudevent(e gjson.Result) event.Event {
	event := c.NewEvent()
	event.Context.SetDataSchema(e.Get("dataschema").String())
	event.SetTime(time.Now())
	event.SetID(e.Get("id").String())
	event.SetSource(e.Get("source").String())
	event.SetType(e.Get("type").String())
	event.SetData(c.ApplicationJSON, e.Get("data"))
	return event
}

func PostHandler(forwarder f.Forwarder, cache *cache.SchemaCache, conf *config.Cloudevents, meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var cloudevents []event.Event
		if c.ContentType() == "application/cloudevents+json" { // Only accept request if content type is set appropriately
			reqBody, _ := ioutil.ReadAll(c.Request.Body)
			result := gjson.ParseBytes(reqBody)
			cloudevent := buildCloudevent(result)
			cloudevents = append(cloudevents, cloudevent)
			validEvents, invalidEvents := bifurcateEvents(cloudevents, cache)
			f.BatchPublishValidAndInvalid(input.CLOUDEVENTS_INPUT, forwarder, validEvents, invalidEvents, meta)
			c.JSON(200, response.Ok)
		} else {
			c.JSON(400, response.InvalidContentType)
		}

	}
	return gin.HandlerFunc(fn)
}

// func BatchPostHandler(forwarder f.Forwarder, cache *cache.SchemaCache, meta *tele.Meta) gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		if c.ContentType() == "application/cloudevents-batch+json" {
// 		} else {
// 			c.JSON(400, response.InvalidContentType)
// 		}
// 	}
// 	return gin.HandlerFunc(fn)
// }
