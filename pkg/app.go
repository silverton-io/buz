package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/env"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
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
	app.config.App.InstanceId = instanceId.String()

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

func (app *App) initializeSnowplowRoutes() {
	log.Info().Msg("initializing snowplow routes")
	log.Info().Msg("initializing health check route")
	app.engine.GET(health.HEALTH_PATH, health.HealthcheckHandler)
	if app.config.Routing.DisableStandardRoutes {
		log.Info().Msg("skipping standard route initialization")
	} else {
		log.Info().Msg("initializing standard routes")
		app.engine.GET(snowplow.DEFAULT_GET_PATH, snowplow.GetHandler(app.forwarder, app.schemaCache))
		app.engine.POST(snowplow.DEFAULT_POST_PATH, snowplow.PostHandler(app.forwarder, app.schemaCache))
		if app.config.Routing.DisableOpenRedirect {
			log.Info().Msg("skipping standard open redirect initialization")
		} else {
			log.Info().Msg("initializing standard open redirect route")
			app.engine.GET(snowplow.DEFAULT_REDIRECT_PATH, snowplow.RedirectHandler(app.forwarder, app.schemaCache))
		}
	}
	log.Info().Msg("initializing custom routes")
	app.engine.GET(app.config.Routing.GetPath, snowplow.GetHandler(app.forwarder, app.schemaCache))
	app.engine.POST(app.config.Routing.PostPath, snowplow.PostHandler(app.forwarder, app.schemaCache))
	if app.config.Routing.DisableOpenRedirect {
		log.Info().Msg("skipping custom open redirect initialization")
	} else {
		log.Info().Msg("initializing custom open redirect route")
		app.engine.GET(app.config.Routing.RedirectPath, snowplow.RedirectHandler(app.forwarder, app.schemaCache))
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
	app.initializeSnowplowRoutes()
	app.serveStaticIfDev()
}

func (app *App) Run() {
	log.Info().Interface("config", app.config).Msg("gosnowplow running with configuration")
	tele.Metry(*app.config)
	app.engine.Run(":" + app.config.App.Port)
}
