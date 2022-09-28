// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package manifold

import (
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/sink"
)

type Manifold interface {
	Initialize(sinks *[]sink.Sink) error
	Distribute(e []envelope.Envelope) error
}
