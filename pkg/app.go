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
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	ce "github.com/silverton-io/gosnowplow/pkg/cloudevents"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/env"
	"github.com/silverton-io/gosnowplow/pkg/generic"
	"github.com/silverton-io/gosnowplow/pkg/health"
	"github.com/silverton-io/gosnowplow/pkg/middleware"
	"github.com/silverton-io/gosnowplow/pkg/sink"
	"github.com/silverton-io/gosnowplow/pkg/snowplow"
	"github.com/silverton-io/gosnowplow/pkg/stats"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"github.com/spf13/viper"
)

type App struct {
	config      *config.Config
	engine      *gin.Engine
	sink        sink.Sink
	schemaCache *cache.SchemaCache
	meta        *tele.Meta
}

var VERSION string

func (a *App) configure() {
	// Set up app logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Load app config from file
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal().Msg("could not read config")
	}
	a.config = &config.Config{}
	viper.Unmarshal(a.config)
	gin.SetMode(a.config.App.Mode)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	a.config.App.Version = VERSION
	instanceId := uuid.New()
	m := tele.Meta{
		Version:       VERSION,
		InstanceId:    instanceId,
		StartTime:     time.Now(),
		TrackerDomain: a.config.App.TrackerDomain,
		CookieDomain:  a.config.Cookie.Domain,
	}
	a.meta = &m
}

func (a *App) initializeSink() {
	log.Info().Msg("initializing sink")
	s, _ := sink.BuildSink(a.config.Sink)
	a.sink = s
}

func (a *App) initializeSchemaCache() {
	log.Info().Msg("initializing schema cache")
	cache := cache.SchemaCache{}
	cache.Initialize(a.config.SchemaCache)
	a.schemaCache = &cache
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
	if a.config.Cookie.Enabled {
		a.engine.Use(middleware.AdvancingCookie(a.config.Cookie))
	}
	a.engine.Use(middleware.Yeet())
	a.engine.Use(middleware.CORS(a.config.Cors))
	a.engine.Use(middleware.JsonAccessLogger())
}

func (a *App) initializeHealthcheckRoutes() {
	log.Info().Msg("initializing health check route")
	a.engine.GET(health.HEALTH_PATH, health.HealthcheckHandler)
}

func (a *App) initializeStatsRoutes() {
	if a.config.App.Stats.Enabled {
		log.Info().Msg("intializing stats route")
		var statsPath string
		if a.config.App.Stats.Endpoint == "" {
			statsPath = stats.STATS_PATH
		} else {
			statsPath = a.config.App.Stats.Endpoint
		}
		a.engine.GET(statsPath, stats.StatsHandler(a.meta))
	}
}

func (a *App) initializeSchemaCachePurgeRoute() {
	if a.config.SchemaCache.Purge.Enabled {
		log.Info().Msg("initializing schema cache purge route")
		a.engine.GET(a.config.SchemaCache.Purge.Path, cache.CachePurgeHandler(a.schemaCache))
	}
}

func (a *App) initializeSnowplowRoutes() {
	if a.config.Snowplow.Enabled {
		log.Info().Msg("initializing snowplow routes")
		if a.config.Snowplow.StandardRoutesEnabled {
			log.Info().Msg("initializing standard routes")
			a.engine.GET(snowplow.DEFAULT_GET_PATH, snowplow.DefaultHandler(a.config.Snowplow, a.meta, a.schemaCache, a.sink))
			a.engine.POST(snowplow.DEFAULT_POST_PATH, snowplow.DefaultHandler(a.config.Snowplow, a.meta, a.schemaCache, a.sink))
			if a.config.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("initializing standard open redirect route")
				a.engine.GET(snowplow.DEFAULT_REDIRECT_PATH, snowplow.RedirectHandler(a.config.Snowplow, a.meta, a.schemaCache, a.sink))
			}
		}
		log.Info().Msg("initializing custom routes")
		a.engine.GET(a.config.Snowplow.GetPath, snowplow.DefaultHandler(a.config.Snowplow, a.meta, a.schemaCache, a.sink))
		a.engine.POST(a.config.Snowplow.PostPath, snowplow.DefaultHandler(a.config.Snowplow, a.meta, a.schemaCache, a.sink))
		if a.config.Snowplow.OpenRedirectsEnabled {
			log.Info().Msg("initializing custom open redirect route")
			a.engine.GET(a.config.Snowplow.RedirectPath, snowplow.RedirectHandler(a.config.Snowplow, a.meta, a.schemaCache, a.sink))
		}
	}
}

func (a *App) initializeGenericRoutes() {
	if a.config.Generic.Enabled {
		log.Info().Msg("initializing generic routes")
		a.engine.POST(a.config.Generic.PostPath, generic.PostHandler(&a.config.Generic, a.meta, a.schemaCache, a.sink))
		a.engine.POST(a.config.Generic.BatchPostPath, generic.BatchPostHandler(&a.config.Generic, a.meta, a.schemaCache, a.sink))
	}
}

func (a *App) initializeCloudeventsRoutes() {
	if a.config.Cloudevents.Enabled {
		log.Info().Msg("initializing cloudevents routes")
		a.engine.POST(a.config.Cloudevents.PostPath, ce.PostHandler(&a.config.Cloudevents, a.meta, a.schemaCache, a.sink))
		a.engine.POST(a.config.Cloudevents.BatchPostPath, ce.BatchPostHandler(&a.config.Cloudevents, a.meta, a.schemaCache, a.sink))
	}
}

func (a *App) serveStaticIfDev() {
	if a.config.App.Env == env.DEV_ENVIRONMENT {
		log.Info().Msg("serving static files")
		a.engine.StaticFile("/", "./static/index.html")     // Serve a local file to make testing events easier
		a.engine.StaticFile("/test", "./static/index.html") // Ditto
	} else {
		log.Info().Msg("not serving static files")
	}
}

func (a *App) Initialize() {
	log.Info().Msg("initializing app")
	a.configure()
	a.initializeSink()
	a.initializeSchemaCache()
	a.initializeRouter()
	a.initializeMiddleware()
	a.initializeHealthcheckRoutes()
	a.initializeStatsRoutes()
	a.initializeSchemaCachePurgeRoute()
	a.initializeSnowplowRoutes()
	a.initializeGenericRoutes()
	a.initializeCloudeventsRoutes()
	a.serveStaticIfDev()
}

func (a *App) Run() {
	log.Info().Interface("config", a.config).Msg("gosnowplow is running!")
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
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Msg("server forced to shutdown")
	}
	tele.Sis(a.meta)
}
