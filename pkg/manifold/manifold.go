package manifold

import (
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/sink"
)

type Manifold interface {
	Initialize(sinks *[]sink.Sink) error
	Distribute(e []envelope.Envelope) error
}
