package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func HandlePost(c *gin.Context) {
	// FIXME! Parse POST Data
	// events := map[string]interface{}{}

	body, _ := ioutil.ReadAll(c.Request.Body)
	data := gjson.GetBytes(body, "data")
	for _, i := range data.Array() {
		fmt.Println(i.Get("e"))
	}
	// err := json.Unmarshal(data, &events)

	c.JSON(200, gin.H{
		"message": "received",
	})
}
