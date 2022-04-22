package handler

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/manifold"
	"github.com/silverton-io/honeypot/pkg/tele"
)

type EventHandlerParams struct {
	Config   *config.Config
	Cache    *cache.SchemaCache
	Manifold *manifold.SimpleManifold
	Meta     *tele.Meta
}
