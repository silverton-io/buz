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
	"github.com/spf13/viper"
)

type App struct {
	config             config.Config
	engine             *gin.Engine
	pubsubClient       *pubsub.Client
	validEventsTopic   *pubsub.Topic
	invalidEventsTopic *pubsub.Topic
}

func (app *App) configure() {
	log.Info().Msg("configuring app")
	viper.SetConfigFile("config.yml")
	viper.ReadInConfig()
	// app.advancingCookieConfig = config.AdvancingCookie{
	// 	CookieName:      "sp-nuid",
	// 	UseSecureCookie: false,
	// 	CookieTtlDays:   365,
	// 	CookiePath:      "/",
	// 	CookieDomain:    "localhost",
	// }
	// app.corsConfig = config.Cors{
	// 	AllowOrigin:      []string{"*"},
	// 	AllowCredentials: true,
	// 	AllowMethods:     []string{"POST", "OPTIONS", "GET"},
	// 	MaxAge:           86400,
	// }
	// app.pubsubConfig = config.Pubsub{
	// 	ProjectName:       "neat-dispatch-338321",
	// 	ValidEventTopic:   "test-valid-topic",
	// 	InvalidEventTopic: "test-invalid-topic",
	// }
	// app.appConfig = config.App{
	// 	Env:                   "dev",
	// 	Port:                  "8080",
	// 	IncludeStandardRoutes: true,
	// SnowplowPostPath:              "com.snowplowanalytics.snowplow/tp2",
	// SnowplowGetPath:               "/i",
	// SnowplowRedirectPath:          "r/tp2",
	// ValidateSnowplowEvents:        false,
	// }
}

func (app *App) initializePubsub() {
	log.Info().Msg("initializing pubsub")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, app.pubsubConfig.ProjectName)
	if err != nil {
		log.Fatal().Msg("could not initialize pubsub client")
	}
	app.pubsubClient = client
	app.validEventsTopic = app.pubsubClient.Topic(app.pubsubConfig.ValidEventTopic)
	app.invalidEventsTopic = app.pubsubClient.Topic(app.pubsubConfig.InvalidEventTopic)
	//defer app.pubsubClient.Close() -> This line results in things blocking forever. Ask MGL why.
}

func (app *App) initializeMiddleware() {
	log.Info().Msg("initializing middleware")
	app.engine.Use(middleware.AdvancingCookie(app.advancingCookieConfig))
	app.engine.Use(middleware.CORS(app.corsConfig))
}

func (app *App) initializeRouter() {
	log.Info().Msg("initializing router")
	app.engine = gin.Default()
	app.engine.RedirectTrailingSlash = false
}

func (app *App) initializeRoutes() {
	log.Info().Msg("initializing routes")
	app.engine.GET(snowplow.DEFAULT_SNOWPLOW_HEALTH_PATH, handler.Healthcheck)
	app.engine.GET(app.appConfig.SnowplowGetPath, handler.SnowplowGet(app.validEventsTopic))
	app.engine.POST(app.appConfig.SnowplowPostPath, handler.SnowplowPost(app.validEventsTopic))
}

func (app *App) serveStaticIfDev() {
	if app.appConfig.Environment == env.DEV_ENVIRONMENT {
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
	app.engine.Run(":" + app.appConfig.Port)
}
