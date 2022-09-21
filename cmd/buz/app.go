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
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/env"
	"github.com/silverton-io/buz/pkg/handler"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/middleware"
	"github.com/silverton-io/buz/pkg/params"
	"github.com/silverton-io/buz/pkg/protocol"
	"github.com/silverton-io/buz/pkg/registry"
	"github.com/silverton-io/buz/pkg/sink"
	"github.com/silverton-io/buz/pkg/stats"
	"github.com/silverton-io/buz/pkg/tele"
	"github.com/spf13/viper"
)

var VERSION string

type App struct {
	config        *config.Config
	engine        *gin.Engine
	registry      *registry.Registry
	manifold      *manifold.SimpleManifold
	sinks         []sink.Sink
	collectorMeta *meta.CollectorMeta
	stats         *stats.ProtocolStats
}

func (a *App) handlerParams() params.Handler {
	params := params.Handler{
		Config:        a.config,
		Registry:      a.registry,
		Manifold:      a.manifold,
		CollectorMeta: a.collectorMeta,
		ProtocolStats: a.stats,
	}
	return params
}

func (a *App) configure() {
	// Set up app logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	gin.SetMode("release")

	// Load app config from file
	conf := os.Getenv(env.BUZ_CONFIG_PATH)
	debug := os.Getenv(env.DEBUG)
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
	if err := viper.Unmarshal(a.config); err != nil {
		log.Fatal().Stack().Err(err).Msg("could not unmarshal config")
	}
	if debug != "" {
		gin.SetMode("debug")
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

func (a *App) initializeRegistry() {
	log.Info().Msg("🟢 initializing schema registry")
	registry := registry.Registry{}
	if err := registry.Initialize(a.config.Registry); err != nil {
		panic(err)
	}
	a.registry = &registry
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
	if err := a.engine.SetTrustedProxies(nil); err != nil {
		panic(err)
	}
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
	log.Info().Msg("🟢 initializing buz route")
	a.engine.GET("/", handler.BuzHandler())
	log.Info().Msg("🟢 initializing health check route")
	a.engine.GET(constants.HEALTH_PATH, handler.HealthcheckHandler)
	log.Info().Msg("🟢 initializing stats route")
	a.engine.GET(constants.STATS_PATH, handler.StatsHandler(a.collectorMeta, a.stats))
	log.Info().Msg("🟢 initializing overview routes")
	a.engine.GET(constants.ROUTE_OVERVIEW_PATH, handler.RouteOverviewHandler(*a.config))
	if a.config.App.EnableConfigRoute {
		log.Info().Msg("🟢 initializing config overview")
		a.engine.GET(constants.CONFIG_OVERVIEW_PATH, handler.ConfigOverviewHandler(*a.config))
	}
}

func (a *App) initializeSchemaCacheRoutes() {
	if a.config.Registry.Purge.Enabled {
		log.Info().Msg("🟢 initializing schema registry cache purge route")
		a.engine.GET(a.config.Registry.Purge.Path, handler.RegistryCachePurgeHandler(a.registry))
	}
	if a.config.Registry.Http.Enabled {
		log.Info().Msg("🟢 initializing schema registry routes")
		a.engine.GET(registry.SCHEMAS_ROUTE+"*"+registry.SCHEMA_PARAM, handler.RegistryGetSchemaHandler(a.registry))
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

func (a *App) initializeSelfDescribingRoutes() {
	if a.config.Inputs.SelfDescribing.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("🟢 initializing generic routes")
		a.engine.POST(a.config.Inputs.SelfDescribing.Path, handler.SelfDescribingHandler(handlerParams))
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
		a.engine.POST(a.config.Inputs.Webhook.Path+"/*"+constants.BUZ_SCHEMA_PARAM, handler.WebhookHandler(handlerParams))
	}
}

func (a *App) initializePixelRoutes() {
	if a.config.Inputs.Pixel.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("🟢 initializing pixel routes")
		a.engine.GET(a.config.Inputs.Pixel.Path, handler.PixelHandler(handlerParams))
		a.engine.GET(a.config.Inputs.Pixel.Path+"/*"+constants.BUZ_SCHEMA_PARAM, handler.PixelHandler(handlerParams))
	}
}

func (a *App) initializeSquawkboxRoutes() {
	if a.config.Squawkbox.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("🟢 initializing squawkbox routes")
		a.engine.POST(constants.SQUAWKBOX_CLOUDEVENTS_PATH, handler.SquawkboxHandler(handlerParams, protocol.CLOUDEVENTS))
		a.engine.POST(constants.SQUAWKBOX_SNOWPLOW_PATH, handler.SquawkboxHandler(handlerParams, protocol.SNOWPLOW))
		a.engine.GET(constants.SQUAWKBOX_SNOWPLOW_PATH, handler.SquawkboxHandler(handlerParams, protocol.SNOWPLOW))
		a.engine.POST(constants.SQUAWKBOX_SELF_DESCRIBING_PATH, handler.SquawkboxHandler(handlerParams, protocol.SELF_DESCRIBING))
		a.engine.GET(constants.SQUAWKBOX_PIXEL_PATH, handler.SquawkboxHandler(handlerParams, protocol.PIXEL))
		a.engine.POST(constants.SQUAWKBOX_WEBHOOK_PATH, handler.SquawkboxHandler(handlerParams, protocol.WEBHOOK))
	}
}

func (a *App) Initialize() {
	log.Info().Msg("🟢 initializing app")
	a.configure()
	a.initializeStats()
	a.initializeSinks()
	a.initializeManifold()
	a.initializeRegistry()
	a.initializeRouter()
	a.initializeMiddleware()
	a.initializeOpsRoutes()
	a.initializeSchemaCacheRoutes()
	a.initializeSnowplowRoutes()
	a.initializeSelfDescribingRoutes()
	a.initializeCloudeventsRoutes()
	a.initializeWebhookRoutes()
	a.initializePixelRoutes()
	a.initializeSquawkboxRoutes()
}

func (a *App) Run() {
	log.Info().Interface("config", a.config).Msg("🐝🐝🐝 buz is running! 🐝🐝🐝")
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
