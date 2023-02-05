// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputselfdescribing

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
)

type SelfDescribingInput struct{}

func (i *SelfDescribingInput) Initialize(engine *gin.Engine, manifold *manifold.Manifold, conf *config.Config, metadata *meta.CollectorMeta) error {
	if conf.Inputs.SelfDescribing.Enabled {
		log.Info().Msg("🟢 initializing generic routes")
		engine.POST(conf.Inputs.SelfDescribing.Path, Handler(*manifold, *conf, metadata))
	}
	return nil
}
