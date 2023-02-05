// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package manifold

import (
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/registry"
	"github.com/silverton-io/buz/pkg/sink"
)

type Manifold interface {
	Initialize(registry *registry.Registry, sinks *[]sink.Sink, conf *config.Config, metadata *meta.CollectorMeta) error
	Distribute(envelopes []envelope.Envelope) error
	GetRegistry() *registry.Registry
	Shutdown() error
}
