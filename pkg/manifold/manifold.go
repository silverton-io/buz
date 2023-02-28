// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package manifold

import (
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/registry"
)

type Manifold interface {
	Initialize(registry *registry.Registry, sinks *[]backendutils.Sink, conf *config.Config, metadata *meta.CollectorMeta) error
	Enqueue(envelopes []envelope.Envelope) error
	GetRegistry() *registry.Registry
	Shutdown() error
}
