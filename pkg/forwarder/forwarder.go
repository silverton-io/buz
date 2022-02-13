package forwarder

import (
	"github.com/silverton-io/gosnowplow/pkg/config"
	"golang.org/x/net/context"
)

type Forwarder interface {
	Initialize(config config.Forwarder)
	PublishValidEvent(ctx context.Context, event interface{})
	PublishInvalidEvent(ctx context.Context, event interface{})
	PublishValidEvents(ctx context.Context, events []interface{})
	PublishInvalidEvents(ctx context.Context, events []interface{})
}
