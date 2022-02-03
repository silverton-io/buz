package main

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/env"
	"github.com/silverton-io/gosnowplow/pkg/handler"
	"github.com/silverton-io/gosnowplow/pkg/middleware"
	"github.com/silverton-io/gosnowplow/pkg/snowplow"
	"github.com/silverton-io/gosnowplow/pkg/util"
	"github.com/spf13/viper"
)

type App struct {
	config             *config.Config
	engine             *gin.Engine
	pubsubClient       *pubsub.Client
	validEventsTopic   *pubsub.Topic
	invalidEventsTopic *pubsub.Topic
}

func (app *App) configure() {
	// Load app config from file
	log.Info().Msg("configuring app")
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Msg("could not read config")
	}
	app.config = &config.Config{}
	viper.Unmarshal(app.config)
	util.PrettyPrint(app.config)

	// Configure gin
	gin.SetMode(app.config.App.Mode)
}

func (app *App) initializePubsub() {
	log.Info().Msg("initializing pubsub")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, app.config.Pubsub.Project)
	if err != nil {
		log.Fatal().Msg("could not initialize pubsub client")
	}
	app.pubsubClient = client
	app.validEventsTopic = app.pubsubClient.Topic(app.config.Pubsub.ValidEventTopic)
	app.invalidEventsTopic = app.pubsubClient.Topic(app.config.Pubsub.InvalidEventTopic)
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
	if app.config.Routing.IncludeStandardRoutes {
		log.Info().Msg("initializing standard routes")
		app.engine.GET(snowplow.DEFAULT_GET_PATH, handler.SnowplowGet(app.validEventsTopic))
		app.engine.POST(snowplow.DEFAULT_POST_PATH, handler.SnowplowPost(app.validEventsTopic))
		// app.engine.GET((snowplow.DEFAULT_REDIRECT_PATH, handler.SnowplowRedirect(app.validEventsTopic)))
	} else {
		log.Info().Msg("skipping standard route initialization")
	}
	app.engine.GET(snowplow.DEFAULT_HEALTH_PATH, handler.Healthcheck)
	app.engine.GET(app.config.Routing.GetPath, handler.SnowplowGet(app.validEventsTopic))
	app.engine.POST(app.config.Routing.PostPath, handler.SnowplowPost(app.validEventsTopic))
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
	app.initializePubsub()
	app.initializeRouter()
	app.initializeMiddleware()
	app.initializeRoutes()
	app.serveStaticIfDev()
}

func (app *App) Run() {
	defer app.pubsubClient.Close()
	app.engine.Run(":" + app.config.App.Port)
}
