package generic

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/tidwall/gjson"
)

func BuildEvent(e gjson.Result, conf config.Generic) GenericEvent {
	var sdContexts []event.SelfDescribingContext
	var sdPayload event.SelfDescribingPayload
	eventContexts := e.Get(conf.Contexts.RootKey).Array()
	if len(eventContexts) == 0 {
		log.Error().Stack().Msg("no contexts found in generic event for key: " + conf.Contexts.RootKey)
	} else {
		for _, context := range eventContexts {
			contextSchema := context.Get(conf.Contexts.SchemaKey).String()
			contextPayload := context.Get(conf.Contexts.DataKey).Value()
			if contextPayload == nil {
				log.Error().Stack().Msg("no contexts payload found in generic event for key: " + conf.Contexts.RootKey + "." + conf.Contexts.DataKey)
				log.Debug().Interface("event", e.Value()).Interface("config", conf).Msg("event format does not match config format")
			} else {
				context := event.SelfDescribingContext{
					Schema: contextSchema,
					Data:   contextPayload.(map[string]interface{}),
				}
				sdContexts = append(sdContexts, context)
			}
		}
	}
	payload := e.Get(conf.Payload.RootKey)
	payloadSchema := payload.Get(conf.Payload.SchemaKey).String()
	payloadData := payload.Get(conf.Payload.DataKey).Value()
	if payloadData == nil {
		log.Error().Stack().Msg("no data payload found in generic event for key: " + conf.Payload.RootKey + "." + conf.Payload.DataKey)
		log.Debug().Interface("event", e.Value()).Interface("config", conf).Msg("event format does not match config format")
	} else {
		sdPayload = event.SelfDescribingPayload{
			Schema: payloadSchema,
			Data:   payloadData.(map[string]interface{}),
		}
	}
	genEvent := GenericEvent{
		Contexts: sdContexts, // FIXME! These contexts are not validated. And should be.
		Payload:  sdPayload,
	}
	return genEvent
}
