package snowplow

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	e "github.com/silverton-io/gosnowplow/pkg/event"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	f "github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/request"
	"github.com/silverton-io/gosnowplow/pkg/response"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"github.com/tidwall/gjson"
)

func bifurcateEvents(events []Event, cache *cache.SchemaCache) (validEvents []interface{}, invalidEvents []interface{}) {
	var vEvents []interface{}
	var invEvents []interface{}
	for _, event := range events {
		isValid, validationError, schema := validateEvent(event, cache)
		setEventMetadataFields(&event, schema)
		if isValid {
			vEvents = append(vEvents, event)
		} else {
			invalidEvent := e.InvalidEvent{
				ValidationError: &validationError,
				Event:           &event,
			}
			invEvents = append(invEvents, invalidEvent)
		}
	}
	return vEvents, invEvents
}

func buildEventsFromRequest(c *gin.Context, conf config.Snowplow) []Event {
	var events []Event
	if c.Request.Method == "POST" {
		body, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs here
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			event := BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), conf)
			events = append(events, event)
		}
	} else {
		mappedParams := request.MapParams(c)
		event := BuildEventFromMappedParams(c, mappedParams, conf)
		events = append(events, event)
	}
	return events
}

func RedirectHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, meta *tele.Meta, conf config.Snowplow) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		events := buildEventsFromRequest(c, conf)
		validEvents, invalidEvents := bifurcateEvents(events, cache)
		f.BatchPublishValidAndInvalid(input.SNOWPLOW_INPUT, forwarder, validEvents, invalidEvents, meta)
		redirectUrl, _ := c.GetQuery("u")
		c.Redirect(302, redirectUrl)
	}
	return gin.HandlerFunc(fn)
}

func DefaultHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, meta *tele.Meta, conf config.Snowplow) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		events := buildEventsFromRequest(c, conf)
		validEvents, invalidEvents := bifurcateEvents(events, cache)
		f.BatchPublishValidAndInvalid(input.SNOWPLOW_INPUT, forwarder, validEvents, invalidEvents, meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
