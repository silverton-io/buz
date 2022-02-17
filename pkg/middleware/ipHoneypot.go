package middleware

import "github.com/gin-gonic/gin"

const RICK string = "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

func IpHoneypot() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
