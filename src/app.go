package main

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type App struct {
	appConfig             AppConfig
	advancingCookieConfig AdvancingCookieConfig
	corsConfig            CORSConfig
	pubsubConfig          PubsubConfig
	engine                *gin.Engine
	pubsubClient          *pubsub.Client
	validEventsTopic      *pubsub.Topic
	invalidEventsTopic    *pubsub.Topic
}

func (app *App) configure() {
	log.Info().Msg("configuring app")
	app.advancingCookieConfig = AdvancingCookieConfig{
		cookieName:      "sp-nuid",
		useSecureCookie: false,
		cookieTtlDays:   365,
		cookiePath:      "/",
		cookieDomain:    "localhost",
	}
	app.corsConfig = CORSConfig{
		allowOrigin:      []string{"*"},
		allowCredentials: true,
		allowMethods:     []string{"POST", "OPTIONS", "GET"},
		maxAge:           86400,
	}
	app.pubsubConfig = PubsubConfig{
		projectName:       "dev-sixrs-data-platform",
		validEventTopic:   "test-valid-topic",
		invalidEventTopic: "test-invalid-topic",
	}
	app.appConfig = AppConfig{
		environment:                   "dev",
		port:                          "8080",
		includeStandardSnowplowRoutes: true,
		snowplowPostPath:              "com.snowplowanalytics.snowplow/tp2",
		snowplowGetPath:               "/i",
		snowplowRedirectPath:          "r/tp2",
		validateSnowplowEvents:        false,
	}
}

func (app *App) initializePubsub() {
	log.Info().Msg("initializing pubsub")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, app.pubsubConfig.projectName)
	if err != nil {
		log.Fatal().Msg("could not initialize pubsub client")
	}
	app.pubsubClient = client
	app.validEventsTopic = app.pubsubClient.Topic(app.pubsubConfig.validEventTopic)
	app.invalidEventsTopic = app.pubsubClient.Topic(app.pubsubConfig.invalidEventTopic)
	//defer app.pubsubClient.Close() -> This line results in things blocking forever. Ask MGL why.
}

func (app *App) initializeMiddleware() {
	log.Info().Msg("initializing middleware")
	app.engine.Use(AdvancingCookieMiddleware(app.advancingCookieConfig))
	app.engine.Use(CORSMiddleware(app.corsConfig))
}

func (app *App) initializeRouter() {
	log.Info().Msg("initializing router")
	app.engine = gin.Default()
	app.engine.RedirectTrailingSlash = false
}

func (app *App) initializeRoutes() {
	log.Info().Msg("initializing routes")
	app.engine.GET(HEALTH_ENDPOINT, HandleHealthcheck)
	app.engine.GET(app.appConfig.snowplowGetPath, HandleGet(app.validEventsTopic))
	app.engine.POST(app.appConfig.snowplowPostPath, HandlePost(app.validEventsTopic))
}

func (app *App) serveStaticIfDev() {
	if app.appConfig.environment == DEV_ENVIRONMENT {
		log.Info().Msg("serving static files")
		// Serve a local file to make testing events easier
		app.engine.StaticFile("/", "./src/static/index.html")
		app.engine.StaticFile("/test/there", "./src/static/index.html")
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
	app.engine.Run(":" + app.appConfig.port)
}
