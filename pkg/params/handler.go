// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

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
