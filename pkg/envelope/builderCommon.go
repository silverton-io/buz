// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/util"
)

func BuildCommonEnvelope(c *gin.Context, conf config.Middleware, m *meta.CollectorMeta) Envelope {
	identity := util.GetIdentityOrFallback(c, conf)
	now := time.Now().UTC()
	envelope := Envelope{
		EventMeta: EventMeta{
			Uuid:      uuid.New(),
			Namespace: constants.UNKNOWN,
		},
		Pipeline: Pipeline{
			Source: Source{
				GeneratedTstamp: &now,
				SentTstamp:      &now,
			},
			Collector: Collector{
				Tstamp:  now,
				Name:    &m.Name,
				Version: &m.Version,
			},
		},
		Device: Device{
			Ip:        c.ClientIP(),
			Id:        identity,
			Useragent: c.Request.UserAgent(),
		},
	}
	return envelope
}
