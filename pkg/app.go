package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/cache"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/env"
	"github.com/silverton-io/gosnowplow/pkg/forwarder"
	"github.com/silverton-io/gosnowplow/pkg/handler"
	"github.com/silverton-io/gosnowplow/pkg/middleware"
	"github.com/silverton-io/gosnowplow/pkg/snowplow"
	"github.com/spf13/viper"
)

type App struct {
	config      *config.Config
	engine      *gin.Engine
	forwarder   *forwarder.PubsubForwarder
	schemaCache *cache.SchemaCache
}

func (app *App) configure() {
	// Load app config from file
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Msg("could not read config")
	}
	app.config = &config.Config{}
	viper.Unmarshal(app.config)
	log.Debug().Interface("config", app.config).Msg("configuring app")
	// Configure gin
	gin.SetMode(app.config.App.Mode)
}

func (app *App) initializeForwarder() {
	log.Info().Msg("initializing forwarder")
	forwarder := forwarder.PubsubForwarder{}
	forwarder.Initialize(app.config.Forwarder)
}

func (app *App) initializeSchemaCache() {
	log.Info().Msg("initializing schema cache")
	cache := cache.SchemaCache{}
	cache.Initialize(app.config.Cache)
	app.schemaCache = &cache
}

func (app *App) initializeRouter() {
	log.Info().Msg("initializing router")
	app.engine = gin.Default()
	app.engine.RedirectTrailingSlash = false
}

func (app *App) initializeMiddleware() {
	log.Info().Msg("initializing middleware")
	app.engine.Use(middleware.AdvancingCookie(app.config.Cookie))
	app.engine.Use(middleware.CORS(app.config.Cors))
}

func (app *App) initializeRoutes() {
	log.Info().Msg("initializing routes")
	log.Info().Msg("initializing health check route")
	app.engine.GET(snowplow.DEFAULT_HEALTH_PATH, handler.Healthcheck)
	if app.config.Routing.DisableStandardRoutes {
		log.Info().Msg("skipping standard route initialization")
	} else {
		log.Info().Msg("initializing standard routes")
		app.engine.GET(snowplow.DEFAULT_GET_PATH, handler.SnowplowGet(app.forwarder))
		app.engine.POST(snowplow.DEFAULT_POST_PATH, handler.SnowplowPost(app.forwarder))
		if app.config.Routing.DisableOpenRedirect {
			log.Info().Msg("skipping standard open redirect initialization")
		} else {
			log.Info().Msg("initializing standard open redirect route")
			app.engine.GET(snowplow.DEFAULT_REDIRECT_PATH, handler.SnowplowRedirect(app.forwarder))
		}
	}
	log.Info().Msg("initializing custom routes")
	app.engine.GET(app.config.Routing.GetPath, handler.SnowplowGet(app.forwarder))
	app.engine.POST(app.config.Routing.PostPath, handler.SnowplowPost(app.forwarder))
	if app.config.Routing.DisableOpenRedirect {
		log.Info().Msg("skipping custom open redirect initialization")
	} else {
		log.Info().Msg("initializing custom open redirect route")
		app.engine.GET(app.config.Routing.RedirectPath, handler.SnowplowRedirect(app.forwarder))
	}
}

func (app *App) serveStaticIfDev() {
	if app.config.App.Env == env.DEV_ENVIRONMENT {
		log.Info().Msg("serving static files")
		// Serve a local file to make testing events easier
		app.engine.StaticFile("/", "./static/index.html")
		app.engine.StaticFile("/test/there", "./static/index.html")
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
	app.initializeRoutes()
	app.serveStaticIfDev()
}

func (app *App) Run() {
	// defer app.forwarder.Client.Close()
	app.engine.Run(":" + app.config.App.Port)
}
