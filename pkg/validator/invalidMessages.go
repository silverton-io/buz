package validator

type InvalidMessage struct {
	Type       string `json:"type"`
	Resolution string `json:"resolution"`
}

var InvalidSchema = InvalidMessage{
	Type:       "invalid schema",
	Resolution: "ensure schema is properly formatted",
}

var InvalidPayload = InvalidMessage{
	Type:       "invalid payload",
	Resolution: "publish a valid payload",
}

var PayloadNotPresent = InvalidMessage{
	Type:       "payload not present",
	Resolution: "publish the event with a payload",
}

var UnknownSnowplowEventType = InvalidMessage{
	Type:       "unknown snowplow event type",
	Resolution: "event type should adhere to the snowplow tracker protocol",
}

var NoSchemaAssociated = InvalidMessage{
	Type:       "no schema associated",
	Resolution: "associate a schema with the event",
}

var NoSchemaInBackend = InvalidMessage{
	Type:       "schema not published to cache backend",
	Resolution: "publish schema to the cache backed",
}
