package params

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/manifold"
	"github.com/silverton-io/honeypot/pkg/tele"
)

type Handler struct {
	Config   *config.Config
	Cache    *cache.SchemaCache
	Manifold *manifold.SimpleManifold
	Meta     *tele.Meta
}
