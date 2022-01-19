package main

import "github.com/gin-gonic/gin"

func HandlePost(c *gin.Context) {
	// FIXME! Parse POST Data
	c.JSON(200, gin.H{
		"message": "received",
	})
}
