package pixel

import "github.com/gin-gonic/gin"

const UNNAMED_PIXEL_ID = "unnamed"

func BuildEvent(c *gin.Context, payload map[string]interface{}) (PixelEvent, error) {
	event := PixelEvent{
		Id:      UNNAMED_PIXEL_ID,
		Payload: payload,
	}
	return event, nil
}
