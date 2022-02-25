package snowplow

import (
	"context"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	e "github.com/silverton-io/gosnowplow/pkg/event"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/request"
	"github.com/silverton-io/gosnowplow/pkg/response"
	"github.com/silverton-io/gosnowplow/pkg/sink"
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

func buildEventsFromRequest(c *gin.Context, s config.Snowplow, t *tele.Meta) []Event {
	var events []Event
	if c.Request.Method == "POST" {
		body, _ := ioutil.ReadAll(c.Request.Body) // FIXME! Handle errs here
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			event := BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), s, t)
			events = append(events, event)
		}
	} else {
		params := request.MapParams(c)
		event := BuildEventFromMappedParams(c, params, s, t)
		events = append(events, event)
	}
	return events
}

func RedirectHandler(sink sink.Sink, cache *cache.SchemaCache, conf config.Snowplow, meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		events := buildEventsFromRequest(c, conf, meta)
		validEvents, invalidEvents := bifurcateEvents(events, cache)
		sink.BatchPublishValidAndInvalid(ctx, input.SNOWPLOW_INPUT, validEvents, invalidEvents, meta)
		redirectUrl, _ := c.GetQuery("u")
		c.Redirect(302, redirectUrl)
	}
	return gin.HandlerFunc(fn)
}

func DefaultHandler(sink sink.Sink, cache *cache.SchemaCache, conf config.Snowplow, meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		ctx := context.Background()
		events := buildEventsFromRequest(c, conf, meta)
		validEvents, invalidEvents := bifurcateEvents(events, cache)
		sink.BatchPublishValidAndInvalid(ctx, input.SNOWPLOW_INPUT, validEvents, invalidEvents, meta)
		c.JSON(200, response.Ok)
	}
	return gin.HandlerFunc(fn)
}
