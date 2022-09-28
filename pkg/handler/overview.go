// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	r "github.com/silverton-io/buz/pkg/registry"
)

type systemPaths struct {
	Health         string `json:"health"`
	Stats          string `json:"stats"`
	RouteOverview  string `json:"routeOverview"`
	ConfigOverview string `json:"configOverview"`
}

type snowplowPaths struct {
	Get      []string `json:"get"`
	Post     []string `json:"post"`
	Redirect []string `json:"redirect"`
}

type inputPaths struct {
	Cloudevents    string        `json:"cloudevents"`
	SelfDescribing string        `json:"selfDescribing"`
	Webhook        string        `json:"webhook"`
	Pixel          string        `json:"pixel"`
	Snowplow       snowplowPaths `json:"snowplow"`
}

type registryPaths struct {
	Base string `json:"base"`
}

type RoutesResponse struct {
	systemPaths   `json:"system"`
	inputPaths    `json:"input"`
	registryPaths `json:"registry"`
}

func RouteOverviewHandler(conf config.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		sp := snowplowPaths{
			Get:      []string{constants.SNOWPLOW_STANDARD_GET_PATH, conf.Snowplow.GetPath},
			Post:     []string{constants.SNOWPLOW_STANDARD_POST_PATH, conf.Snowplow.PostPath},
			Redirect: []string{constants.SNOWPLOW_STANDARD_REDIRECT_PATH, conf.Snowplow.RedirectPath},
		}
		resp := RoutesResponse{
			systemPaths{
				Health:         constants.HEALTH_PATH,
				Stats:          constants.STATS_PATH,
				RouteOverview:  constants.ROUTE_OVERVIEW_PATH,
				ConfigOverview: constants.CONFIG_OVERVIEW_PATH,
			},
			inputPaths{
				Snowplow:       sp,
				Cloudevents:    conf.Cloudevents.Path,
				SelfDescribing: conf.SelfDescribing.Path,
				Webhook:        conf.Webhook.Path,
				Pixel:          conf.Pixel.Path,
			},
			registryPaths{
				Base: r.SCHEMAS_ROUTE,
			},
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
