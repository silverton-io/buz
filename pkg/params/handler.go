// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package params

import (
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/registry"
	"github.com/silverton-io/buz/pkg/stats"
)

type Handler struct {
	Config        *config.Config
	Registry      *registry.Registry
	Manifold      manifold.Manifold
	CollectorMeta *meta.CollectorMeta
	ProtocolStats *stats.ProtocolStats
}
