package handler

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/sink"
	"github.com/silverton-io/honeypot/pkg/tele"
)

type EventHandlerParams struct {
	Config *config.Config
	Cache  *cache.SchemaCache
	Sink   sink.Sink
	Meta   *tele.Meta
}
