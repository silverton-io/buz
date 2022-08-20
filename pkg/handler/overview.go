// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
)

type RoutesResponse struct {
	HealthPath                   string `json:"healthPath"`
	StatsPath                    string `json:"statsPath"`
	RouteOverviewPath            string `json:"routeOverviewPath"`
	ConfigOverviewPath           string `json:"configOverviewPath"`
	CloudeventsPath              string `json:"cloudeventsPath"`
	GenericPath                  string `json:"genericPath"`
	WebhookPath                  string `json:"webhookPath"`
	PixelPath                    string `json:"pixelPath"`
	SnowplowStandardGetPath      string `json:"snowplowStandardGetPath"`
	SnowplowGetPath              string `json:"snowplowGetPath"`
	SnowplowStandardPostPath     string `json:"snowplowStandardPostPath"`
	SnowplowPostPath             string `json:"snowplowPostPath"`
	SnowplowStandardRedirectPath string `json:"snowplowStandardRedirectPath"`
	SnowplowRedirectPath         string `json:"snowplowRedirectPath"`
}

func RouteOverviewHandler(conf config.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		resp := RoutesResponse{
			HealthPath:                   constants.HEALTH_PATH,
			StatsPath:                    constants.STATS_PATH,
			RouteOverviewPath:            constants.ROUTE_OVERVIEW_PATH,
			ConfigOverviewPath:           constants.CONFIG_OVERVIEW_PATH,
			CloudeventsPath:              conf.Cloudevents.Path,
			GenericPath:                  conf.Generic.Path,
			WebhookPath:                  conf.Webhook.Path,
			PixelPath:                    conf.Pixel.Path,
			SnowplowStandardGetPath:      constants.SNOWPLOW_STANDARD_GET_PATH,
			SnowplowGetPath:              conf.Snowplow.GetPath,
			SnowplowStandardPostPath:     constants.SNOWPLOW_STANDARD_POST_PATH,
			SnowplowPostPath:             conf.Snowplow.PostPath,
			SnowplowStandardRedirectPath: constants.SNOWPLOW_STANDARD_REDIRECT_PATH,
			SnowplowRedirectPath:         conf.Snowplow.RedirectPath,
		}
		c.JSON(200, resp)
	}
	return gin.HandlerFunc(fn)
}

func ConfigOverviewHandler(conf config.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(200, conf)
	}
	return gin.HandlerFunc(fn)
}
