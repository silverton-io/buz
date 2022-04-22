package handler

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/manifold"
)

type EventHandlerParams struct {
	Config   *config.Config
	Cache    *cache.SchemaCache
	Manifold *manifold.SimpleManifold
}
