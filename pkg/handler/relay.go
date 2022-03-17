package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/silverton-io/honeypot/pkg/validator"
)

func RelayHandler(p EventHandlerParams) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var envelopes []event.Envelope
		ctx := context.Background()
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		json.Unmarshal(reqBody, &envelopes)
		util.Pprint(envelopes)
		validEnvelopes, invalidEnvelopes := validator.Bifurcate(envelopes)
	}
	return gin.HandlerFunc(fn)
}
