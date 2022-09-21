// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package protocol

const (
	SNOWPLOW        string = "snowplow"
	SELF_DESCRIBING string = "selfDescribing"
	CLOUDEVENTS     string = "cloudevents"
	WEBHOOK         string = "webhook"
	PIXEL           string = "pixel"
)

func GetIntputProtocols() []string {
	return []string{SNOWPLOW, SELF_DESCRIBING, CLOUDEVENTS, WEBHOOK, PIXEL}
}
