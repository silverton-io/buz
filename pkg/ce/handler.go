package ce

import (
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	f "github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/input"
	"github.com/silverton-io/gosnowplow/pkg/tele"
)

func bifurcateEvents(events event.Event, cache *cache.SchemaCache) (validEvents []event.Event, invalidEvents []event.Event) {
	var vEvents []event.Event
	var invEvents []event.Event
}

func PostHandler(forwarder forwarder.Forwarder, cache *cache.SchemaCache, meta *tele.Meta) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		validEvents, invalidEvents := bifurcateEvents(events, cache)
		f.BatchPublishValidAndInvalid(input.CLOUDEVENTS_INPUT, forwarder, validEvents, invalidEvents, meta)

	}
	return gin.HandlerFunc(fn)
}

func BatchPostHandler(forwarder f.Forwarder, cache *cache.SchemaCache, meta *tele.Meta)
