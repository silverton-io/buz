package manifold

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/sink"
	"github.com/silverton-io/honeypot/pkg/tele"
)

type Manifold struct {
	envelopeChan          *chan envelope.Envelope
	bufferRecordThreshold int
	bufferByteThreshold   int
	bufferTimeThreshold   int
	buffersLastPurged     time.Time
	sinks                 *[]sink.Sink
	validEnvelopes        []envelope.Envelope
	invalidEnvelopes      []envelope.Envelope
}

func (m *Manifold) initialize(conf config.Manifold, sinks *[]sink.Sink) {
	c := make(chan envelope.Envelope, 10000) // FIXME! Configurable buffer size
	m.envelopeChan = &c
	m.bufferRecordThreshold = conf.BufferRecordThreshold
	m.bufferByteThreshold = conf.BufferByteThreshold
	m.bufferTimeThreshold = conf.BufferTimeThreshold
	m.sinks = sinks
	m.buffersLastPurged = time.Now()
}

func (m Manifold) Enqueue(envelopes []envelope.Envelope) {
	for _, e := range envelopes {
		log.Debug().Msg("enqueing envelope")
		*m.envelopeChan <- e
		// FIXME! Add durable queue option
	}
}

func (m Manifold) validCount() int {
	return len(m.validEnvelopes)
}

func (m Manifold) invalidCount() int {
	return len(m.invalidEnvelopes)
}

func (m Manifold) buffersFull() bool {
	if m.validCount() >= m.bufferRecordThreshold || m.invalidCount() >= m.bufferRecordThreshold {
		return true
	}
	return false
}

func (m *Manifold) PurgeBuffersToSinks(ctx context.Context) {
	// FIXME - what happens when one sink entirely succeeds and one fails (partially or entirely)? Will need better handling of this.
	for _, sink := range *m.sinks {
		m.purgeBuffersToSink(ctx, &sink)
	}
	m.invalidEnvelopes = nil
	m.validEnvelopes = nil
	m.buffersLastPurged = time.Now()
}

func (m *Manifold) PurgeBuffersToSinksIfFull(ctx context.Context) {
	if m.buffersFull() {
		m.PurgeBuffersToSinks(ctx)
	}
}

func (m *Manifold) purgeBuffersToSink(ctx context.Context, s *sink.Sink) {
	sink := *s
	log.Info().Interface("envelopeCount", m.validCount()).Interface("sinkId", sink.Id()).Interface("sinkName", sink.Name()).Msg("sinking valid envelopes")
	sink.BatchPublishValid(ctx, m.validEnvelopes) // FIXME! What happens when sink fails? Should buffer somewhere durably.
	log.Info().Interface("envelopeCount", m.invalidCount()).Interface("sinkId", sink.Id()).Interface("sinkName", sink.Name()).Msg("sinking invalid envelopes")
	sink.BatchPublishInvalid(ctx, m.invalidEnvelopes) // FIXME! What happens when sink fails? Should buffer somewhere durably.
}

func (m *Manifold) AppendValidEnvelope(e envelope.Envelope) {
	v := append(m.validEnvelopes, e)
	m.validEnvelopes = v
}

func (m *Manifold) AppendInvalidEnvelope(e envelope.Envelope) {
	i := append(m.invalidEnvelopes, e)
	m.invalidEnvelopes = i
}

func (m *Manifold) Run(meta *tele.Meta, shutdown *chan bool) {
	log.Debug().Msg("running manifold")
	go func() {
		ctx := context.Background()
		for {
			select {
			case <-*shutdown:
				log.Info().Msg("manifold is shutting down")
				m.PurgeBuffersToSinks(ctx)
				meta.BufferPurgeStats.Increment()
				return
			case e := <-*m.envelopeChan:
				if *e.IsValid {
					log.Debug().Msg("appending valid envelope to buffer...")
					m.AppendValidEnvelope(e)
					meta.ProtocolStats.IncrementValid(e.EventProtocol, e.EventSchema, 1)
				} else {
					log.Debug().Msg("appending invalid envelope to buffer...")
					m.AppendInvalidEnvelope(e)
					meta.ProtocolStats.IncrementInvalid(e.EventProtocol, e.EventSchema, 1)
				}
				m.PurgeBuffersToSinksIfFull(ctx)
				meta.BufferPurgeStats.Increment()
			}
			// FIXME! Have a ticker on a bufferTimeThreshold and send to a channel
		}
	}()
}

func BuildManifold(conf config.Manifold, sinks *[]sink.Sink) (manifold *Manifold, err error) {
	log.Debug().Msg("building manifold")
	m := Manifold{}
	m.initialize(conf, sinks)
	return &m, nil
}
