// Copyright (c) 2023 Silverton Data, Inconf.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package pixel

import (
	"encoding/base64"
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

const PX string = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8Xw8AAoMBgDTD2qgAAAAASUVORK5CYII="

type PixelInput struct{}

func (i *PixelInput) Initialize(engine *gin.Engine, manifold *manifold.Manifold, conf *config.Config, metadata *meta.CollectorMeta) error {
	pixelGroup := engine.Group(conf.Inputs.Pixel.Path)
	if conf.Inputs.Pixel.Enabled {
		log.Info().Msg("🟢 initializing pixel input")
		pixelGroup.GET("", i.Handler(*manifold, *conf, metadata))
		pixelGroup.GET("*"+constants.BUZ_SCHEMA_PARAM, i.Handler(*manifold, *conf, metadata))
	}
	if conf.Squawkbox.Enabled {
		log.Info().Msg("🟢 initializing pixel input squawkbox")
		pixelGroup.GET("squawkbox/pixel", i.SquawkboxHandler(*manifold, *conf, metadata))
	}
	return nil
}

func (i *PixelInput) Handler(m manifold.Manifold, conf config.Config, metadata *meta.CollectorMeta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := i.EnvelopeBuilder(c, &conf, metadata)
		err := m.Enqueue(envelopes)
		if err != nil {
			c.Header("Retry-After", response.RETRY_AFTER_60)
			c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
		} else {
			b, _ := base64.StdEncoding.DecodeString(PX)
			c.Data(http.StatusOK, "image/png", b)
		}
	}
	return gin.HandlerFunc(fn)
}

func (i *PixelInput) SquawkboxHandler(m manifold.Manifold, conf config.Config, metadata *meta.CollectorMeta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := i.EnvelopeBuilder(c, &conf, metadata)
		c.JSON(http.StatusOK, envelopes)
	}
	return gin.HandlerFunc(fn)
}

func (i *PixelInput) EnvelopeBuilder(c *gin.Context, conf *config.Config, metadata *meta.CollectorMeta) []envelope.Envelope {
	return buildEnvelopesFromRequest(c, conf, metadata)
}
