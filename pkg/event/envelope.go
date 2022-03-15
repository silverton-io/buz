package event

import "time"

type Envelope struct {
	Protocol          string                 `json:"protocol"`
	EventVendor       string                 `json:"eventVendor"`
	EventName         string                 `json:"eventName"`
	EventVersion      string                 `json:"eventVersion"`
	EventSchemaFormat string                 `json:"eventSchemaFormat"`
	Tstamp            time.Time              `json:"tstamp"`
	Ip                string                 `json:"ip"`
	Event             map[string]interface{} `json:"event"`
}
