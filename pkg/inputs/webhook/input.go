// Copyright (c) 2023 Silverton Data, Inconf.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputwebhook

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
)

type WebhookInput struct{}

func (i *WebhookInput) Initialize(engine *gin.Engine, manifold *manifold.Manifold, conf *config.Config, metadata *meta.CollectorMeta) error {
	if conf.Inputs.Webhook.Enabled {
		log.Info().Msg("🟢 initializing webhook routes")
		engine.POST(conf.Inputs.Webhook.Path, Handler(*manifold, *conf, metadata))
		engine.POST(conf.Inputs.Webhook.Path+"/*"+constants.BUZ_SCHEMA_PARAM, Handler(*manifold, *conf, metadata))
	}
	return nil
}
