package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.RedirectTrailingSlash = false
	if Config.environment == DEV_ENVIRONMENT {
		// Serve a local file to make testing events easier
		router.StaticFile("/", "./src/static/index.html")
		router.StaticFile("/test/there", "./src/static/index.html")
	}
	router.Use(CORSMiddleware())
	router.Use(AdvancingCookieMiddleware())
	router.GET(HEALTH_ENDPOINT, HandleHealthcheck)
	router.GET(Config.snowplowGetPath, HandleGet)
	router.GET(Config.snowplowRedirectPath, HandleRedirect)
	router.POST(Config.snowplowPostPath, HandlePost)
	router.Run(":" + Config.port)
}
