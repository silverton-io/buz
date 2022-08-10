// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package envelope

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/constants"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/util"
)

func buildCommonEnvelope(c *gin.Context, conf config.Middleware, m *meta.CollectorMeta) Envelope {
	identity := util.GetIdentityOrFallback(c, conf)
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
			Id:        identity,
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
