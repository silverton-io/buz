// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package params

import (
	"github.com/silverton-io/buz/pkg/cache"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/stats"
)

type Handler struct {
	Config        *config.Config
	Cache         *cache.SchemaCache
	Manifold      *manifold.SimpleManifold
	CollectorMeta *meta.CollectorMeta
	ProtocolStats *stats.ProtocolStats
}
