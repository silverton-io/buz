package sink

import (
	"context"
	"net/url"
	"path"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/request"
)

const (
	PUBNUB_PUBLISH_URL   string = "ps.pndsn.com/publish"
	PUBNUB_HTTP_PROTOCOL string = "https"
)

type PubnubSink struct {
	validChannel   string
	invalidChannel string
	pubKey         string
	subKey         string
	store          int
	callback       string
	id             uuid.UUID
}

func (s *PubnubSink) Initialize(conf config.Sink) {
	log.Debug().Msg("initializing pubnub sink")
	s.validChannel = conf.ValidChannel
	s.invalidChannel = conf.InvalidChannel
	s.pubKey = conf.PubnubPubKey
	s.subKey = conf.PubnubSubKey
	u := uuid.New()
	s.id = u
}

func (s *PubnubSink) buildPublishUrl(channel string) *url.URL {
	p := path.Join(PUBNUB_PUBLISH_URL, s.pubKey, s.subKey, "1", channel, "0")
	p = "https://" + p + "?uuid=" + s.id.String()
	u, err := url.Parse(p)
	if err != nil {
		log.Error().Err(err).Msg("could not parse url")
	}
	return u
}

func (s *PubnubSink) batchPublish(ctx context.Context, channel string, envelopes []envelope.Envelope) {
	u := s.buildPublishUrl(channel)
	_, err := request.PostEnvelopes(*u, envelopes)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not post envelopes")
	}
}

func (s *PubnubSink) BatchPublishValid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.validChannel, envelopes)
}

func (s *PubnubSink) BatchPublishInvalid(ctx context.Context, envelopes []envelope.Envelope) {
	s.batchPublish(ctx, s.invalidChannel, envelopes)
}

func (s *PubnubSink) Close() {
	log.Debug().Msg("closing pubnub sink")
}
