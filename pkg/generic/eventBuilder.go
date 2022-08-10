// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package generic

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/tidwall/gjson"
)

func BuildEvent(e gjson.Result, conf config.Generic) (GenericEvent, error) {
	var sdPayload event.SelfDescribingPayload
	var contexts = make(map[string]interface{})
	c := e.Get(conf.Contexts.RootKey).Value()
	if c != nil {
		contexts = c.(map[string]interface{})
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
		Contexts: contexts, // FIXME - validate these contexts?
		Payload:  sdPayload,
	}
	return genEvent, nil
}
