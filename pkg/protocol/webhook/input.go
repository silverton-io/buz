// Copyright (c) 2023 Silverton Data, Inconf.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package webhook

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/response"
)

type WebhookInput struct{}

func (i *WebhookInput) Initialize(routerGroup *gin.RouterGroup, manifold *manifold.Manifold, conf *config.Config, metadata *meta.CollectorMeta) error {
	if conf.Inputs.Webhook.Enabled {
		log.Info().Msg("🟢 initializing webhook input")
		routerGroup.POST(conf.Inputs.Webhook.Path, i.Handler(*manifold, *conf, metadata))
		routerGroup.POST(conf.Inputs.Webhook.Path+"/*"+constants.BUZ_SCHEMA_PARAM, i.Handler(*manifold, *conf, metadata))
	}
	if conf.Squawkbox.Enabled {
		log.Info().Msg("🟢 initializing webhook input squawkbox")
		routerGroup.POST("/squawkbox/webhook", i.SquawkboxHandler(*manifold, *conf, metadata))
	}
	return nil
}

func (i *WebhookInput) Handler(m manifold.Manifold, conf config.Config, metadata *meta.CollectorMeta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			envelopes := i.EnvelopeBuilder(c, &conf, metadata)
			err := m.Enqueue(envelopes)
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

func (i *WebhookInput) SquawkboxHandler(m manifold.Manifold, conf config.Config, metadata *meta.CollectorMeta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := i.EnvelopeBuilder(c, &conf, metadata)
		c.JSON(http.StatusOK, envelopes)
	}
	return gin.HandlerFunc(fn)
}

func (i *WebhookInput) EnvelopeBuilder(c *gin.Context, conf *config.Config, metadata *meta.CollectorMeta) []envelope.Envelope {
	return buildEnvelopesFromRequest(c, conf, metadata)
}
