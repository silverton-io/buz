// Copyright (c) 2022 Silverton Data, Inconf.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package snowplow

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/middleware"
	"github.com/silverton-io/buz/pkg/response"
)

type SnowplowInput struct{}

func (i *SnowplowInput) Initialize(engine *gin.Engine, manifold *manifold.Manifold, conf *config.Config, metadata *meta.CollectorMeta) error {
	identityMiddleware := middleware.Identity(conf.Identity)
	if conf.Inputs.Snowplow.Enabled {
		log.Info().Msg("🟢 initializing snowplow routes")
		if conf.Inputs.Snowplow.StandardRoutesEnabled {
			log.Info().Msg("🟢 initializing standard snowplow routes")
			engine.GET(constants.SNOWPLOW_STANDARD_GET_PATH, identityMiddleware, i.Handler(*manifold, *conf, metadata))
			engine.POST(constants.SNOWPLOW_STANDARD_POST_PATH, identityMiddleware, i.Handler(*manifold, *conf, metadata))
			if conf.Inputs.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("🟢 initializing standard open redirect route")
				engine.GET(constants.SNOWPLOW_STANDARD_REDIRECT_PATH, identityMiddleware, i.Handler(*manifold, *conf, metadata))
			}
		}
		log.Info().Msg("🟢 initializing custom snowplow routes")
		engine.GET(conf.Inputs.Snowplow.GetPath, identityMiddleware, i.Handler(*manifold, *conf, metadata))
		engine.POST(conf.Inputs.Snowplow.PostPath, identityMiddleware, i.Handler(*manifold, *conf, metadata))
		if conf.Inputs.Snowplow.OpenRedirectsEnabled {
			log.Info().Msg("🟢 initializing custom open redirect route")
			engine.GET(conf.Inputs.Snowplow.RedirectPath, identityMiddleware, i.Handler(*manifold, *conf, metadata))
		}
	}
	return nil
}

func (i *SnowplowInput) Handler(m manifold.Manifold, conf config.Config, metadata *meta.CollectorMeta) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := BuildEnvelopesFromRequest(c, &conf, metadata)
		err := m.Distribute(envelopes)
		if err != nil {
			c.Header("Retry-After", response.RETRY_AFTER_60)
			c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
		} else {
			c.JSON(http.StatusOK, response.Ok)
		}
		if c.Request.Method == http.MethodGet {
			redirectUrl, _ := c.GetQuery("u")
			if redirectUrl != "" && conf.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("🟢 redirecting to " + redirectUrl)
				c.Redirect(http.StatusFound, redirectUrl)
			}
		}
	}
	return gin.HandlerFunc(fn)
}
