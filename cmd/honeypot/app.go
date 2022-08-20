// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/cache"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/handler"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/middleware"
	"github.com/silverton-io/buz/pkg/params"
	"github.com/silverton-io/buz/pkg/protocol"
	"github.com/silverton-io/buz/pkg/sink"
	"github.com/silverton-io/buz/pkg/stats"
	"github.com/silverton-io/buz/pkg/tele"
	"github.com/spf13/viper"
)

var VERSION string

type App struct {
	config        *config.Config
	engine        *gin.Engine
	schemaCache   *cache.SchemaCache
	manifold      *manifold.SimpleManifold
	sinks         []sink.Sink
	collectorMeta *meta.CollectorMeta
	stats         *stats.ProtocolStats
}

func (a *App) handlerParams() params.Handler {
	params := params.Handler{
		Config:        a.config,
		Cache:         a.schemaCache,
		Manifold:      a.manifold,
		CollectorMeta: a.collectorMeta,
		ProtocolStats: a.stats,
	}
	return params
}

func (a *App) configure() {
	// Set up app logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Load app config from file
	conf := os.Getenv(env.buz_CONFIG_PATH)
	if conf == "" {
		conf = "config.yml"
	}
	log.Info().Msg("🟢 loading config from " + conf)
	viper.SetConfigFile(conf)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not read config")
	}
	a.config = &config.Config{}
	viper.Unmarshal(a.config)
	gin.SetMode(a.config.App.Mode)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	a.config.App.Version = VERSION
	meta := meta.BuildCollectorMeta(VERSION, a.config)
	a.collectorMeta = meta
}

func (a *App) initializeStats() {
	log.Info().Msg("🟢 initializing stats")
	ps := stats.ProtocolStats{}
	ps.Build()
	a.stats = &ps
}

func (a *App) initializeSchemaCache() {
	log.Info().Msg("🟢 initializing schema cache")
	cache := cache.SchemaCache{}
	cache.Initialize(a.config.SchemaCache)
	a.schemaCache = &cache
}

func (a *App) initializeSinks() {
	log.Info().Msg("🟢 initializing sinks")
	sinks, err := sink.BuildAndInitializeSinks(a.config.Sinks)
	if err != nil {
		log.Fatal().Err(err).Msg("could not build and init sinks")
	}
	a.sinks = sinks
}

func (a *App) initializeManifold() {
	log.Info().Msg("🟢 initializing manifold")
	manifold := manifold.SimpleManifold{}
	err := manifold.Initialize(&a.sinks)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not build manifold")
	}
	a.manifold = &manifold
}

func (a *App) initializeRouter() {
	log.Info().Msg("🟢 initializing router")
	a.engine = gin.New()
	a.engine.SetTrustedProxies(nil)
	a.engine.RedirectTrailingSlash = false
}

func (a *App) initializeMiddleware() {
	log.Info().Msg("🟢 initializing middleware")
	a.engine.Use(gin.Recovery())
	if a.config.Middleware.Timeout.Enabled {
		log.Info().Msg("🟢 initializing request timeout middleware")
		a.engine.Use(middleware.Timeout(a.config.Middleware.Timeout))
	}
	if a.config.Middleware.RateLimiter.Enabled {
		log.Info().Msg("🟢 initializing rate limiter middleware")
		limiter := middleware.BuildRateLimiter(a.config.Middleware.RateLimiter)
		limiterMiddleware := middleware.BuildRateLimiterMiddleware(limiter)
		a.engine.Use(limiterMiddleware)
	}
	if a.config.Middleware.Cors.Enabled {
		log.Info().Msg("🟢 initializing cors middleware")
		a.engine.Use(middleware.CORS(a.config.Middleware.Cors))
	}
	if a.config.Middleware.RequestLogger.Enabled {
		log.Info().Msg("🟢 initializing request logger middleware")
		a.engine.Use(middleware.RequestLogger())
	}
	if a.config.Middleware.Yeet.Enabled {
		log.Info().Msg("🟢 initializing yeet middleware")
		a.engine.Use(middleware.Yeet())
	}
	log.Info().Msg("🟢 initializing identity middleware")
	a.engine.Use(middleware.Identity(a.config.Identity))
}

func (a *App) initializeOpsRoutes() {
	log.Info().Msg("🟢 initializing health check route")
	a.engine.GET(constants.HEALTH_PATH, handler.HealthcheckHandler)
	log.Info().Msg("🟢 intializing stats route")
	a.engine.GET(constants.STATS_PATH, handler.StatsHandler(a.collectorMeta, a.stats))
	log.Info().Msg("🟢 initializing overview routes")
	a.engine.GET(constants.ROUTE_OVERVIEW_PATH, handler.RouteOverviewHandler(*a.config))
	if a.config.App.EnableConfigRoute {
		log.Info().Msg("🟢 initializing config overview")
		a.engine.GET(constants.CONFIG_OVERVIEW_PATH, handler.ConfigOverviewHandler(*a.config))
	}
}

func (a *App) initializeSchemaCacheRoutes() {
	if a.config.SchemaCache.Purge.Enabled {
		log.Info().Msg("🟢 initializing schema cache purge route")
		a.engine.GET(a.config.SchemaCache.Purge.Path, handler.CachePurgeHandler(a.schemaCache))
	}
	if a.config.SchemaCache.SchemaDirectory.Enabled {
		log.Info().Msg("🟢 initializing schema cache index and getter routes")
		a.engine.GET(cache.SCHEMA_CACHE_ROOT_ROUTE, handler.CacheIndexHandler(a.schemaCache))
		a.engine.GET(cache.SCHEMA_CACHE_ROOT_ROUTE+"/*"+cache.SCHEMA_ROUTE_PARAM, handler.CacheGetHandler(a.schemaCache))
	}
}

func (a *App) initializeSnowplowRoutes() {
	if a.config.Inputs.Snowplow.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("🟢 initializing snowplow routes")
		if a.config.Inputs.Snowplow.StandardRoutesEnabled {
			log.Info().Msg("🟢 initializing standard snowplow routes")
			a.engine.GET(constants.SNOWPLOW_STANDARD_GET_PATH, handler.SnowplowHandler(handlerParams))
			a.engine.POST(constants.SNOWPLOW_STANDARD_POST_PATH, handler.SnowplowHandler(handlerParams))
			if a.config.Inputs.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("🟢 initializing standard open redirect route")
				a.engine.GET(constants.SNOWPLOW_STANDARD_REDIRECT_PATH, handler.SnowplowHandler(handlerParams))
			}
		}
		log.Info().Msg("🟢 initializing custom snowplow routes")
		a.engine.GET(a.config.Inputs.Snowplow.GetPath, handler.SnowplowHandler(handlerParams))
		a.engine.POST(a.config.Inputs.Snowplow.PostPath, handler.SnowplowHandler(handlerParams))
		if a.config.Inputs.Snowplow.OpenRedirectsEnabled {
			log.Info().Msg("🟢 initializing custom open redirect route")
			a.engine.GET(a.config.Inputs.Snowplow.RedirectPath, handler.SnowplowHandler(handlerParams))
		}
	}
}

func (a *App) initializeGenericRoutes() {
	if a.config.Inputs.Generic.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("🟢 initializing generic routes")
		a.engine.POST(a.config.Inputs.Generic.Path, handler.GenericHandler(handlerParams))
	}
}

func (a *App) initializeCloudeventsRoutes() {
	if a.config.Inputs.Cloudevents.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("🟢 initializing cloudevents routes")
		a.engine.POST(a.config.Inputs.Cloudevents.Path, handler.CloudeventsHandler(handlerParams))
	}
}

func (a *App) initializeWebhookRoutes() {
	if a.config.Inputs.Webhook.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("🟢 initializing webhook routes")
		a.engine.POST(a.config.Inputs.Webhook.Path, handler.WebhookHandler(handlerParams))
		a.engine.POST(a.config.Inputs.Webhook.Path+"/*"+constants.buz_SCHEMA_PARAM, handler.WebhookHandler(handlerParams))
	}
}

func (a *App) initializePixelRoutes() {
	if a.config.Inputs.Pixel.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("🟢 initializing pixel routes")
		a.engine.GET(a.config.Inputs.Pixel.Path, handler.PixelHandler(handlerParams))
		a.engine.GET(a.config.Inputs.Pixel.Path+"/*"+constants.buz_SCHEMA_PARAM, handler.PixelHandler(handlerParams))
	}
}

func (a *App) initializeSquawkboxRoutes() {
	if a.config.Squawkbox.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("🟢 initializing squawkbox routes")
		a.engine.POST(a.config.Squawkbox.CloudeventsPath, handler.SquawkboxHandler(handlerParams, protocol.CLOUDEVENTS))
		a.engine.POST(a.config.Squawkbox.GenericPath, handler.SquawkboxHandler(handlerParams, protocol.GENERIC))
		a.engine.POST(a.config.Squawkbox.SnowplowPath, handler.SquawkboxHandler(handlerParams, protocol.SNOWPLOW))
		a.engine.GET(a.config.Squawkbox.SnowplowPath, handler.SquawkboxHandler(handlerParams, protocol.SNOWPLOW))
	}
}

func (a *App) Initialize() {
	log.Info().Msg("🟢 initializing app")
	a.configure()
	a.initializeStats()
	a.initializeSinks()
	a.initializeManifold()
	a.initializeSchemaCache()
	a.initializeRouter()
	a.initializeMiddleware()
	a.initializeOpsRoutes()
	a.initializeSchemaCacheRoutes()
	a.initializeSnowplowRoutes()
	a.initializeGenericRoutes()
	a.initializeCloudeventsRoutes()
	a.initializeWebhookRoutes()
	a.initializePixelRoutes()
	a.initializeSquawkboxRoutes()
}

func (a *App) Run() {
	log.Info().Interface("config", a.config).Msg("🍯🍯🍯 buz is running! 🍯🍯🍯")
	tele.Metry(a.config, a.collectorMeta)
	srv := &http.Server{
		Addr:    ":" + a.config.App.Port,
		Handler: a.engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Info().Msgf("🟢 server shut down")
		}
	}()
	// Safe shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("🟢 shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Stack().Err(err).Msg("server forced to shutdown")
	}
	tele.Sis(a.collectorMeta)
}
