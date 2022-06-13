package envelope

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/constants"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/meta"
)

func buildCommonEnvelope(c *gin.Context, m *meta.CollectorMeta) Envelope {
	envelope := Envelope{
		EventMeta: EventMeta{
			Uuid:      uuid.New(),
			Namespace: constants.UNKNOWN,
		},
		Pipeline: Pipeline{
			Source: Source{
				GeneratedTstamp: time.Now().UTC(),
				SentTstamp:      time.Now().UTC(),
			},
			Collector: Collector{
				Tstamp:  time.Now().UTC(),
				Name:    &m.Name,
				Version: &m.Version,
			},
			Relay: Relay{
				Relayed: false,
			},
		},
		Device: Device{
			Ip:        c.ClientIP(),
			Useragent: c.Request.UserAgent(),
		},
		User:       User{},
		Session:    Session{},
		Web:        Web{},
		Validation: Validation{},
		Contexts:   event.Contexts{},
	}
	return envelope
}
