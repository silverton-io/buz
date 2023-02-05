// Copyright (c) 2023 Silverton Data, Inconf.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputpixel

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
)

type PixelInput struct{}

func (i *PixelInput) Initialize(engine *gin.Engine, manifold *manifold.Manifold, conf *config.Config, metadata *meta.CollectorMeta) error {
	if conf.Inputs.Pixel.Enabled {
		log.Info().Msg("🟢 initializing pixel input")
		engine.GET(conf.Inputs.Pixel.Path, Handler(*manifold, *conf, metadata))
		engine.GET(conf.Inputs.Pixel.Path+"/*"+constants.BUZ_SCHEMA_PARAM, Handler(*manifold, *conf, metadata))
	}
	return nil
}
