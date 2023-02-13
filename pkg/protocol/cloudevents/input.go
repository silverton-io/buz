// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package cloudevents

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/response"
)

type CloudeventsInput struct{}

func (i *CloudeventsInput) Initialize(engine *gin.Engine, manifold *manifold.Manifold, conf *config.Config, metadata *meta.CollectorMeta) error {
	if conf.Inputs.Cloudevents.Enabled {
		log.Info().Msg("🟢 initializing cloudevents input")
		engine.POST(conf.Inputs.Cloudevents.Path, i.Handler(*manifold, *conf, metadata))
	}
	if conf.Squawkbox.Enabled {
		log.Info().Msg("🟢 initializing cloudevents input squawkbox")
		engine.POST("/squawkbox/cloudevents", i.SquawkboxHandler(*manifold, *conf, metadata))
	}
	return nil
}

func (i *CloudeventsInput) Handler(m manifold.Manifold, conf config.Config, metadata *meta.CollectorMeta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.ContentType() == "application/cloudevents+json" || c.ContentType() == "application/cloudevents-batch+json" {
			envelopes := i.EnvelopeBuilder(c, &conf, metadata)
			err := m.Distribute(envelopes)
			if err != nil {
				c.Header("Retry-After", response.RETRY_AFTER_60)
				c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
			} else {
				c.JSON(http.StatusOK, response.Ok)
			}
		} else {
			c.JSON(http.StatusBadRequest, response.InvalidContentType)
		}
	}
	return gin.HandlerFunc(fn)

}

func (i *CloudeventsInput) SquawkboxHandler(m manifold.Manifold, conf config.Config, metadata *meta.CollectorMeta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := i.EnvelopeBuilder(c, &conf, metadata)
		c.JSON(http.StatusOK, envelopes)
	}
	return gin.HandlerFunc(fn)
}

func (i *CloudeventsInput) EnvelopeBuilder(c *gin.Context, conf *config.Config, metadata *meta.CollectorMeta) []envelope.Envelope {
	return buildEnvelopesFromRequest(c, conf, metadata)
}
