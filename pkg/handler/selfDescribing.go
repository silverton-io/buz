package handler

// import (
// 	"context"
// 	"encoding/json"
// 	"io/ioutil"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rs/zerolog/log"
// 	"github.com/silverton-io/honeypot/pkg/cache"
// 	"github.com/silverton-io/honeypot/pkg/config"
// 	e "github.com/silverton-io/honeypot/pkg/event"
// 	"github.com/silverton-io/honeypot/pkg/protocol"
// 	"github.com/silverton-io/honeypot/pkg/response"
// 	"github.com/silverton-io/honeypot/pkg/sink"
// 	"github.com/silverton-io/honeypot/pkg/tele"
// 	"github.com/tidwall/gjson"
// )

// func PostHandler(conf *config.Generic, meta *tele.Meta, cache *cache.SchemaCache, sink sink.Sink) gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		ctx := context.Background()
// 		reqBody, _ := ioutil.ReadAll(c.Request.Body)
// 		var events []gjson.Result
// 		event := gjson.ParseBytes(reqBody)
// 		events = append(events, event)
// 		validEvents, invalidEvents := bifurcateEvents(events, cache, conf)
// 		sink.BatchPublishValidAndInvalid(ctx, protocol.GENERIC, validEvents, invalidEvents, meta)
// 		c.JSON(200, response.Ok)
// 	}
// 	return gin.HandlerFunc(fn)
// }

// func BatchPostHandler(conf *config.Generic, meta *tele.Meta, cache *cache.SchemaCache, sink sink.Sink) gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		ctx := context.Background()
// 		reqBody, _ := ioutil.ReadAll(c.Request.Body)
// 		var rawEvents []interface{}
// 		var events []gjson.Result
// 		err := json.Unmarshal(reqBody, &rawEvents)
// 		if err != nil {
// 			log.Error().Stack().Err(err).Msg("error when unmarshaling request body")
// 			// TODO! Decide whether or not to return something bad here
// 		}
// 		for _, rawEvent := range rawEvents {
// 			marshaledEvent, err := json.Marshal(rawEvent)
// 			if err != nil {
// 				log.Error().Stack().Err(err).Msg("error when marshaling event")
// 				// TODO! Decide whether or not to return something bad here
// 			} else {
// 				event := gjson.ParseBytes(marshaledEvent)
// 				events = append(events, event)
// 			}
// 		}
// 		validEvents, invalidEvents := bifurcateEvents(events, cache, conf)
// 		sink.BatchPublishValidAndInvalid(ctx, protocol.GENERIC, validEvents, invalidEvents, meta)
// 		c.JSON(200, response.Ok)
// 	}
// 	return gin.HandlerFunc(fn)
// }
