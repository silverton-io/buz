package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func mapParams(c *gin.Context) map[string]interface{} {
	mappedParams := make(map[string]interface{})
	params := c.Request.URL.Query()
	for k, v := range params {
		mappedParams[k] = v[0]
	}
	return mappedParams
}

func buildEventFromMappedParams(e map[string]interface{}) Event {
	body, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err)
	}
	shortenedEvent := ShortenedEvent{}
	json.Unmarshal(body, &shortenedEvent)
	return Event(shortenedEvent)
}

func HandleRedirect(c *gin.Context) {
	mappedParams := mapParams(c)
	se := buildEventFromMappedParams(mappedParams)
	fmt.Printf("\n%+v\n", se)
	// e, _ := json.Marshal(se)
	redirectUrl, _ := c.GetQuery("u")
	c.Redirect(302, redirectUrl)
}

func HandleGet(c *gin.Context) {
	// Parse query parameters to map[string]interface{}
	mappedParams := mapParams(c)
	se := buildEventFromMappedParams(mappedParams)
	fmt.Printf("\n%+v\n", se)
	// e, _ := json.Marshal(se)
	c.JSON(200, AllOk)
}

// func HandlePost(c *gin.Context) {
// 	body, err := ioutil.ReadAll(c.Request.Body)
// 	if err != nil {
// 		fmt.Println("FIXME! HANDLE ERRORS WHEN READING REQUEST BODY")
// 	}

// 	payloadData := gjson.GetBytes(body, "data")
// 	for _, e := range payloadData.Array() {
// 		// event := buildEvent(c, e)
// 		// fmt.Printf("\n%+v\n", event)
// 		fmt.Println(e)
// 	}
// 	c.JSON(200, AllOk)
// }
