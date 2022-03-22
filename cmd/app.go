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
	"github.com/silverton-io/honeypot/pkg/cache"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/env"
	"github.com/silverton-io/honeypot/pkg/handler"
	"github.com/silverton-io/honeypot/pkg/manifold"
	"github.com/silverton-io/honeypot/pkg/middleware"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/sink"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/silverton-io/honeypot/pkg/webhook"
	"github.com/spf13/viper"
)

var VERSION string

type App struct {
	config      *config.Config
	engine      *gin.Engine
	schemaCache *cache.SchemaCache
	manifold    *manifold.Manifold
	sink        sink.Sink
	meta        *tele.Meta
}

func (a *App) handlerParams() handler.EventHandlerParams {
	params := handler.EventHandlerParams{
		Config:   a.config,
		Cache:    a.schemaCache,
		Manifold: a.manifold,
	}
	return params
}

func (a *App) configure() {
	// Set up app logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Load app config from file
	conf := os.Getenv(env.HONEYPOT_CONFIG_PATH)
	if conf == "" {
		conf = "config.yml"
	}
	log.Info().Msg("loading config from " + conf)
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
	meta := tele.BuildMeta(VERSION, a.config)
	a.meta = meta
}

func (a *App) initializeSchemaCache() {
	log.Info().Msg("initializing schema cache")
	cache := cache.SchemaCache{}
	cache.Initialize(a.config.SchemaCache)
	a.schemaCache = &cache
}

func (a *App) initializeSinks() {
	log.Info().Msg("initializing sinks")
	s, _ := sink.BuildSink(a.config.Sink) // FIXME! What happens if the sink creation throws an err?
	sink.InitializeSink(a.config.Sink, s)
	a.sink = s
}

func (a *App) initializeManifold() {
	log.Info().Msg("initializing manifold")
	m, err := manifold.BuildManifold(a.config.Manifold, &a.sink)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not build manifold")
	}
	a.manifold = m
}

func (a *App) initializeRouter() {
	log.Info().Msg("initializing router")
	a.engine = gin.New()
	a.engine.SetTrustedProxies(nil)
	a.engine.RedirectTrailingSlash = false
}

func (a *App) initializeMiddleware() {
	log.Info().Msg("initializing middleware")
	a.engine.Use(gin.Recovery())
	if a.config.Middleware.Timeout.Enabled {
		log.Info().Msg("initializing request timeout middleware")
		a.engine.Use(middleware.Timeout(a.config.Middleware.Timeout))
	}
	if a.config.Middleware.RateLimiter.Enabled {
		log.Info().Msg("initializing rate limiter middleware")
		limiter := middleware.BuildRateLimiter(a.config.Middleware.RateLimiter)
		limiterMiddleware := middleware.BuildRateLimiterMiddleware(limiter)
		a.engine.Use(limiterMiddleware)
	}
	if a.config.Middleware.Cookie.Enabled {
		log.Info().Msg("initializing advancing cookie middleware")
		a.engine.Use(middleware.AdvancingCookie(a.config.Cookie))
	}
	if a.config.Middleware.Cors.Enabled {
		log.Info().Msg("initializing cors middleware")
		a.engine.Use(middleware.CORS(a.config.Middleware.Cors))
	}
	if a.config.Middleware.RequestLogger.Enabled {
		log.Info().Msg("initializing request logger middleware")
		a.engine.Use(middleware.RequestLogger())
	}
	if a.config.Middleware.Yeet.Enabled {
		log.Info().Msg("initializing yeet middleware")
		a.engine.Use(middleware.Yeet())
	}
}

func (a *App) initializeHealthcheckRoutes() {
	if a.config.App.Health.Enabled {
		log.Info().Msg("initializing health check route")
		var healthPath string
		if a.config.App.Health.Path == "" {
			healthPath = "/health"
		} else {
			healthPath = a.config.App.Health.Path
		}
		a.engine.GET(healthPath, handler.HealthcheckHandler)
	}
}

func (a *App) initializeStatsRoutes() {
	if a.config.App.Stats.Enabled {
		log.Info().Msg("intializing stats route")
		var statsPath string
		if a.config.App.Stats.Path == "" {
			statsPath = "/stats"
		} else {
			statsPath = a.config.App.Stats.Path
		}
		a.engine.GET(statsPath, handler.StatsHandler(a.meta))
	}
}

func (a *App) initializeSchemaCacheRoutes() {
	if a.config.SchemaCache.Purge.Enabled {
		log.Info().Msg("initializing schema cache purge route")
		a.engine.GET(a.config.SchemaCache.Purge.Path, handler.CachePurgeHandler(a.schemaCache))
	}
	if a.config.SchemaCache.SchemaEndpoints.Enabled {
		log.Info().Msg("initializing schema cache index and getter routes")
		a.engine.GET(cache.SCHEMA_CACHE_ROOT_ROUTE, handler.CacheIndexHandler(a.schemaCache))
		a.engine.GET(cache.SCHEMA_CACHE_ROOT_ROUTE+"/*"+cache.SCHEMA_ROUTE_PARAM, handler.CacheGetHandler(a.schemaCache))
	}
}

func (a *App) initializeSnowplowRoutes() {
	if a.config.Inputs.Snowplow.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("initializing snowplow routes")
		if a.config.Inputs.Snowplow.StandardRoutesEnabled {
			log.Info().Msg("initializing standard routes")
			a.engine.GET(snowplow.DEFAULT_GET_PATH, handler.SnowplowHandler(handlerParams))
			a.engine.POST(snowplow.DEFAULT_POST_PATH, handler.SnowplowHandler(handlerParams))
			if a.config.Inputs.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("initializing standard open redirect route")
				a.engine.GET(snowplow.DEFAULT_REDIRECT_PATH, handler.SnowplowHandler(handlerParams))
			}
		}
		log.Info().Msg("initializing custom routes")
		a.engine.GET(a.config.Inputs.Snowplow.GetPath, handler.SnowplowHandler(handlerParams))
		a.engine.POST(a.config.Inputs.Snowplow.PostPath, handler.SnowplowHandler(handlerParams))
		if a.config.Inputs.Snowplow.OpenRedirectsEnabled {
			log.Info().Msg("initializing custom open redirect route")
			a.engine.GET(a.config.Inputs.Snowplow.RedirectPath, handler.SnowplowHandler(handlerParams))
		}
	}
}

func (a *App) initializeGenericRoutes() {
	if a.config.Inputs.Generic.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("initializing generic routes")
		a.engine.POST(a.config.Inputs.Generic.PostPath, handler.GenericHandler(handlerParams))
		a.engine.POST(a.config.Inputs.Generic.BatchPostPath, handler.GenericHandler(handlerParams))
	}
}

func (a *App) initializeCloudeventsRoutes() {
	if a.config.Inputs.Cloudevents.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("initializing cloudevents routes")
		a.engine.POST(a.config.Inputs.Cloudevents.PostPath, handler.CloudeventsHandler(handlerParams))
		a.engine.POST(a.config.Inputs.Cloudevents.BatchPostPath, handler.CloudeventsHandler(handlerParams))
	}
}

func (a *App) initializeWebhookRoutes() {
	if a.config.Inputs.Webhook.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("initializing webhook routes")
		a.engine.POST(a.config.Inputs.Webhook.Path, handler.WebhookHandler(handlerParams))
		a.engine.POST(a.config.Inputs.Webhook.Path+"/*"+webhook.WEBHOOK_ID_PARAM, handler.WebhookHandler(handlerParams))
	}
}

func (a *App) initializeRelayRoute() {
	if a.config.Inputs.Relay.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("initializing relay route")
		a.engine.POST(a.config.Inputs.Relay.Path, handler.RelayHandler(handlerParams))
	}
}

func (a *App) initializeSquawkboxRoutes() {
	if a.config.Squawkbox.Enabled {
		handlerParams := a.handlerParams()
		log.Info().Msg("initializing squawkbox routes")
		a.engine.POST(a.config.Squawkbox.CloudeventsPath, handler.SquawkboxHandler(handlerParams, protocol.CLOUDEVENTS))
		a.engine.POST(a.config.Squawkbox.GenericPath, handler.SquawkboxHandler(handlerParams, protocol.GENERIC))
		a.engine.POST(a.config.Squawkbox.SnowplowPath, handler.SquawkboxHandler(handlerParams, protocol.SNOWPLOW))
		a.engine.GET(a.config.Squawkbox.SnowplowPath, handler.SquawkboxHandler(handlerParams, protocol.SNOWPLOW))
	}
}

func (a *App) serveStaticIfDev() {
	if a.config.App.Env == env.DEV_ENVIRONMENT {
		log.Info().Msg("serving static files")
		a.engine.StaticFile("/", "./site/index.html")     // Serve a local file to make testing events easier
		a.engine.StaticFile("/test", "./site/index.html") // Ditto
	} else {
		log.Info().Msg("not serving static files")
	}
}

func (a *App) Initialize() {
	log.Info().Msg("initializing app")
	a.configure()
	a.initializeSinks()
	a.initializeManifold()
	a.initializeSchemaCache()
	a.initializeRouter()
	a.initializeMiddleware()
	a.initializeHealthcheckRoutes()
	a.initializeStatsRoutes()
	a.initializeSchemaCacheRoutes()
	a.initializeSnowplowRoutes()
	a.initializeGenericRoutes()
	a.initializeCloudeventsRoutes()
	a.initializeWebhookRoutes()
	a.initializeSquawkboxRoutes()
	a.initializeRelayRoute()
	a.serveStaticIfDev()
}

func (a *App) Run() {
	log.Info().Interface("config", a.config).Msg("🍯🍯🍯 honeypot is running! 🍯🍯🍯")
	tele.Metry(a.config, a.meta)
	srv := &http.Server{
		Addr:    ":" + a.config.App.Port,
		Handler: a.engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Info().Msgf("server shut down")
		}
	}()
	manifold.Run(a.manifold, a.meta)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Stack().Err(err).Msg("server forced to shutdown")
	}
	tele.Sis(a.meta)
}
