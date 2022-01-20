package main

import (
	"github.com/gin-gonic/gin"
)

func HandleGet(c *gin.Context) {
	// ctx := context.Background()
	// event := eventFromRequest(c)
	// eventBytes, _ := json.Marshal(event)
	// result := PubsubTopic.Publish(ctx, &pubsub.Message{Data: eventBytes})
	// id, err := result.Get(ctx) // FIXME! This is blocking.
	// if err != nil {
	// 	fmt.Println(err) // FIXME!
	// }
	// fmt.Println("Published message. ID: ", id)
	c.JSON(200, gin.H{
		"message": "received",
	})
}

func HandleRedirect(c *gin.Context) {
	// event := eventFromRequest(c)
	// fmt.Printf("%+v\n", event)
	redirectUrl, _ := c.GetQuery("u")
	c.Redirect(302, redirectUrl)
}
