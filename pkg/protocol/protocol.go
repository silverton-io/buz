package protocol

const (
	SNOWPLOW    string = "snowplow"
	GENERIC     string = "generic"
	CLOUDEVENTS string = "cloudevents"
	WEBHOOK     string = "webhook"
	PIXEL       string = "pixel"
)

func GetIntputProtocols() []string {
	return []string{SNOWPLOW, GENERIC, CLOUDEVENTS, WEBHOOK, PIXEL}
}
