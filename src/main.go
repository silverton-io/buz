package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	buildPubsubClient()

	router := gin.Default()
	router.RedirectTrailingSlash = false
	if Config.environment == "dev" {
		// Serve a local file to make testing events easier
		router.StaticFile("/", "./src/static/index.html")
	}
	router.Use(CORSMiddleware())
	router.Use(AdvancingCookieMiddleware())

	// Write
	router.GET(Config.getPath, HandleGet)
	router.GET(Config.redirectPath, HandleRedirect)
	router.POST(Config.postPath, HandlePost)

	router.Run(":" + Config.port)
}
