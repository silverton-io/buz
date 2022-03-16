package generic

import (
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/tidwall/gjson"
)

func BuildEvent(e gjson.Result, conf config.Generic) GenericEvent {

	// contexts := e.Get(conf.Contexts.RootKey)
	// FIXME! For context in contexts coerce to sd context
	payload := e.Get(conf.Payload.RootKey)
	payloadSchema := payload.Get(conf.Payload.SchemaKey).String()
	payloadData := payload.Get(conf.Payload.DataKey).Value().(map[string]interface{})
	sdPayload := event.SelfDescribingPayload{
		Schema: payloadSchema,
		Data:   payloadData,
	}
	genEvent := GenericEvent{
		Payload: sdPayload,
	}
	return genEvent
}
