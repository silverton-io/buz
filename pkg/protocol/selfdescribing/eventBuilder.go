// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package selfdescribing

import (
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/tidwall/gjson"
)

type GenericEvent envelope.SelfDescribingEvent

func buildEvent(e gjson.Result, conf config.SelfDescribing) (GenericEvent, error) {
	var sdPayload envelope.SelfDescribingPayload
	c := e.Get(conf.Contexts.RootKey).Value()
	var contexts map[string]interface{}
	if c != nil {
		contexts = c.(map[string]interface{})
	}
	payload := e.Get(conf.Payload.RootKey)
	payloadSchema := payload.Get(conf.Payload.SchemaKey).String()
	payloadData := payload.Get(conf.Payload.DataKey).Value()
	if payloadData == nil {
		log.Error().Stack().Msg("🔴 no data payload found in generic event for key: " + conf.Payload.RootKey + "." + conf.Payload.DataKey)
		log.Debug().Interface("event", e.Value()).Interface("config", conf).Msg("🟡 event format does not match config format")
	} else {
		sdPayload = envelope.SelfDescribingPayload{
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
