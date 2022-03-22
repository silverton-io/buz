package manifold

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/sink"
)

type Manifold struct {
	envelopeChan      *chan envelope.Envelope
	bufferRecordCount int
	sink              *sink.Sink
	lastPurged        time.Time
}

func (m *Manifold) initialize(conf config.Manifold, sink *sink.Sink) {
	c := make(chan envelope.Envelope, 10000) // FIXME! Configurable buffer size
	m.envelopeChan = &c
	m.bufferRecordCount = conf.BufferRecordThreshold
	m.sink = sink
	m.lastPurged = time.Now()
}

func (m Manifold) Enqueue(envelopes []envelope.Envelope) {
	for _, e := range envelopes {
		log.Debug().Msg("enqueing envelope")
		*m.envelopeChan <- e
	}
}

func BuildManifold(conf config.Manifold, sink *sink.Sink) (manifold *Manifold, err error) {
	log.Debug().Msg("building manifold")
	m := Manifold{}
	m.initialize(conf, sink)
	return &m, nil
}

func Run(m *Manifold) {
	log.Debug().Msg("running manifold")
	go func() {
		ctx := context.Background()
		var invalidEnvelopes []envelope.Envelope
		var validEnvelopes []envelope.Envelope
		sink := *m.sink
		for {
			e := <-*m.envelopeChan
			if *e.IsValid {
				log.Debug().Msg("appending valid envelope to buffer...")
				validEnvelopes = append(validEnvelopes, e)
			} else {
				log.Debug().Msg("appending invalid envelope to buffer...")
				invalidEnvelopes = append(invalidEnvelopes, e)
			}
			if len(validEnvelopes) >= m.bufferRecordCount || len(invalidEnvelopes) >= m.bufferRecordCount {
				log.Debug().Msg("purging envelope buffers")
				sink.BatchPublishValid(ctx, validEnvelopes)
				sink.BatchPublishInvalid(ctx, invalidEnvelopes)
				m.lastPurged = time.Now()
				invalidEnvelopes = nil
				validEnvelopes = nil
			}
		}
	}()
}
