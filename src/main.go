package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	buildPubsubClients()

	router := gin.Default()
	router.RedirectTrailingSlash = false
	if Config.environment == "dev" {
		// Serve a local file to make testing events easier
		router.StaticFile("/", "./src/static/index.html")
	}
	router.Use(CORSMiddleware())
	router.Use(AdvancingCookieMiddleware())

	// Write
	router.GET(Config.snowplowGetPath, HandleGet)
	router.GET(Config.snowplowRedirectPath, HandleRedirect)
	// router.POST(Config.snowplowPostPath, HandlePost)

	router.Run(":" + Config.port)
}
