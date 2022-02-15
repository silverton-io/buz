package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/env"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/generic"
	"github.com/silverton-io/gosnowplow/pkg/health"
	"github.com/silverton-io/gosnowplow/pkg/middleware"
	"github.com/silverton-io/gosnowplow/pkg/snowplow"
	"github.com/silverton-io/gosnowplow/pkg/tele"
	"github.com/spf13/viper"
)

type App struct {
	config      *config.Config
	engine      *gin.Engine
	forwarder   forwarder.Forwarder
	schemaCache *cache.SchemaCache
	meta        *tele.Meta
}

var VERSION string

func (app *App) configure() {
	// Set up app logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Load app config from file
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal().Msg("could not read config")
	}
	app.config = &config.Config{}
	viper.Unmarshal(app.config)
	gin.SetMode(app.config.App.Mode)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// Set version from ldflag
	app.config.App.Version = VERSION
	// Generate and set an instance id
	instanceId := uuid.New()
	m := tele.Meta{
		Version:    VERSION,
		InstanceId: instanceId,
		StartTime:  time.Now(),
		Domain:     app.config.Cookie.Domain,
	}
	app.meta = &m
}

func (app *App) initializeForwarder() {
	log.Info().Msg("initializing forwarder")
	forwarder, _ := forwarder.BuildForwarder(app.config.Forwarder)
	app.forwarder = forwarder
}

func (app *App) initializeSchemaCache() {
	log.Info().Msg("initializing schema cache")
	cache := cache.SchemaCache{}
	cache.Initialize(app.config.SchemaCache)
	app.schemaCache = &cache
}

func (app *App) initializeRouter() {
	log.Info().Msg("initializing router")
	app.engine = gin.New()
	app.engine.RedirectTrailingSlash = false
}

func (app *App) initializeMiddleware() {
	log.Info().Msg("initializing middleware")
	app.engine.Use(gin.Recovery())
	app.engine.Use(middleware.AdvancingCookie(app.config.Cookie))
	app.engine.Use(middleware.CORS(app.config.Cors))
	app.engine.Use(middleware.JsonAccessLogger())
}

func (app *App) initializeHealthcheck() {
	log.Info().Msg("initializing health check route")
	app.engine.GET(health.HEALTH_PATH, health.HealthcheckHandler)
}

func (app *App) initializeSnowplowRoutes() {
	if app.config.Snowplow.Enabled {
		log.Info().Msg("initializing snowplow routes")
		if app.config.Snowplow.StandardRoutesEnabled {
			log.Info().Msg("initializing standard routes")
			app.engine.GET(snowplow.DEFAULT_GET_PATH, snowplow.GetHandler(app.forwarder, app.schemaCache))
			app.engine.POST(snowplow.DEFAULT_POST_PATH, snowplow.PostHandler(app.forwarder, app.schemaCache))
			if app.config.Snowplow.OpenRedirectsEnabled {
				log.Info().Msg("initializing standard open redirect route")
				app.engine.GET(snowplow.DEFAULT_REDIRECT_PATH, snowplow.RedirectHandler(app.forwarder, app.schemaCache))
			}
		}
		log.Info().Msg("initializing custom routes")
		app.engine.GET(app.config.Snowplow.GetPath, snowplow.GetHandler(app.forwarder, app.schemaCache))
		app.engine.POST(app.config.Snowplow.PostPath, snowplow.PostHandler(app.forwarder, app.schemaCache))
		if app.config.Snowplow.OpenRedirectsEnabled {
			log.Info().Msg("initializing custom open redirect route")
			app.engine.GET(app.config.Snowplow.RedirectPath, snowplow.RedirectHandler(app.forwarder, app.schemaCache))
		}
	}
}

func (app *App) initializeGenericRoutes() {
	if app.config.Generic.Enabled {
		log.Info().Msg("initializing generic routes")
		app.engine.POST(app.config.Generic.PostPath, generic.PostHandler(app.forwarder, app.schemaCache, app.config.Generic))
		app.engine.POST(app.config.Generic.BatchPostPath, generic.BatchPostHandler(app.forwarder, app.schemaCache, app.config.Generic))
	}
}

func (app *App) serveStaticIfDev() {
	if app.config.App.Env == env.DEV_ENVIRONMENT {
		log.Info().Msg("serving static files")
		app.engine.StaticFile("/", "./static/index.html")           // Serve a local file to make testing events easier
		app.engine.StaticFile("/test/there", "./static/index.html") // Ditto
	} else {
		log.Info().Msg("not serving static files")
	}
}

func (app *App) Initialize() {
	log.Info().Msg("initializing app")
	app.configure()
	app.initializeForwarder()
	app.initializeSchemaCache()
	app.initializeRouter()
	app.initializeMiddleware()
	app.initializeHealthcheck()
	app.initializeSnowplowRoutes()
	app.initializeGenericRoutes()
	app.serveStaticIfDev()
}

func (app *App) Run() {
	log.Info().Interface("config", app.config).Msg("gosnowplow running with configuration")
	tele.Metry(app.config, app.meta)
	app.engine.Run(":" + app.config.App.Port)
}
