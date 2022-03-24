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
	ShutdownChan          *chan bool
	envelopeChan          *chan envelope.Envelope
	bufferRecordThreshold int
	bufferByteThreshold   int
	bufferTimeThreshold   int
	sinks                 *[]sink.Sink
	lastPurged            time.Time
}

func (m *Manifold) initialize(conf config.Manifold, sinks *[]sink.Sink) {
	c := make(chan envelope.Envelope, 10000) // FIXME! Configurable buffer size
	sDown := make(chan bool, 1)
	m.ShutdownChan = &sDown
	m.envelopeChan = &c
	m.bufferRecordThreshold = conf.BufferRecordThreshold
	m.bufferByteThreshold = conf.BufferByteThreshold
	m.bufferTimeThreshold = conf.BufferTimeThreshold
	m.sinks = sinks
	m.lastPurged = time.Now()
}

func (m Manifold) Enqueue(envelopes []envelope.Envelope) {
	for _, e := range envelopes {
		log.Debug().Msg("enqueing envelope")
		*m.envelopeChan <- e
		// FIXME! Add durability option
	}
}

func BuildManifold(conf config.Manifold, sinks *[]sink.Sink) (manifold *Manifold, err error) {
	log.Debug().Msg("building manifold")
	m := Manifold{}
	m.initialize(conf, sinks)
	return &m, nil
}

func Run(manifold *Manifold, meta *tele.Meta) {
	log.Debug().Msg("running manifold")
	go func() {
		ctx := context.Background()
		var invalidEnvelopes []envelope.Envelope
		var validEnvelopes []envelope.Envelope
		sinks := *manifold.sinks
		for {
			e := <-*manifold.envelopeChan
			if *e.IsValid {
				log.Debug().Msg("appending valid envelope to buffer...")
				validEnvelopes = append(validEnvelopes, e)
				meta.ProtocolStats.IncrementValid(e.EventProtocol, e.EventSchema, 1)
			} else {
				log.Debug().Msg("appending invalid envelope to buffer...")
				invalidEnvelopes = append(invalidEnvelopes, e)
				meta.ProtocolStats.IncrementInvalid(e.EventProtocol, e.EventSchema, 1)
			}
			vEnvelopeCount := len(validEnvelopes)
			invEnvelopeCount := len(invalidEnvelopes)
			if vEnvelopeCount >= manifold.bufferRecordThreshold || invEnvelopeCount >= manifold.bufferRecordThreshold { // FIXME! Break out buffer purge
				log.Debug().Msg("purging envelope buffers")
				for _, sink := range sinks {
					log.Info().Interface("envelopeCount", vEnvelopeCount).Interface("sinkId", sink.Id()).Interface("sinkName", sink.Name()).Msg("sinking valid envelopes")
					sink.BatchPublishValid(ctx, validEnvelopes) // FIXME! What happens when sink fails? Should buffer somewhere durably.
					log.Info().Interface("envelopeCount", invEnvelopeCount).Interface("sinkId", sink.Id()).Interface("sinkName", sink.Name()).Msg("sinking invalid envelopes")
					sink.BatchPublishInvalid(ctx, invalidEnvelopes) // FIXME! What happens when sink fails? Should buffer somewhere durably.
				}
				meta.BufferPurgeStats.Increment()
				manifold.lastPurged = time.Now()
				invalidEnvelopes = nil
				validEnvelopes = nil
			}
		}
	}()
}
