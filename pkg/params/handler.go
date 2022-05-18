package params

import (
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/manifold"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/stats"
)

type Handler struct {
	Config        *config.Config
	Cache         *cache.SchemaCache
	Manifold      *manifold.SimpleManifold
	CollectorMeta *meta.CollectorMeta
	ProtocolStats *stats.ProtocolStats
}
