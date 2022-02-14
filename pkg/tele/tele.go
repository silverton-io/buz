package tele

import (
	"time"

	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/snowplow"
	"github.com/silverton-io/gosnowplow/pkg/util"
)

const (
	DEFAULT_HOST string = ""
)

type configSnapshot struct {
	Version    string        `json:"version"`
	InstanceId string        `json:"instanceId"`
	Domain     string        `json:"domain"`
	Time       time.Time     `json:"time"`
	Config     config.Config `json:"config"`
}

type beat struct {
	Version    string    `json:"version"`
	InstanceId string    `json:"instanceId"`
	Domain     string    `json:"domain"`
	Time       time.Time `json:"time"`
}

func heartbeat(t time.Ticker, c config.Config) {
	for _ = range t.C {
		b := beat{
			Version:    c.App.Version,
			InstanceId: c.App.InstanceId,
			Domain:     c.Cookie.Domain,
			Time:       time.Now(),
		}
		data := util.StructToMap(b)
		sd := snowplow.SelfDescribingPayload{
			Schema: "com.silverton.io.tele/heartbeat/jsonschema/1-0-0",
			Data:   data,
		}
		util.PrettyPrint(sd)
	}
}

func Metry(c config.Config) {
	if c.Tele.Enable {
		snapshot := configSnapshot{
			Version:    c.App.Version,
			InstanceId: c.App.InstanceId,
			Domain:     c.Cookie.Domain,
			Time:       time.Now(),
			Config:     c,
		}
		data := util.StructToMap(snapshot)
		sd := snowplow.SelfDescribingPayload{
			Schema: "com.silverton.io.tele/snapshot/jsonschema/1-0-0",
			Data:   data,
		}
		util.PrettyPrint(sd)
		ticker := time.NewTicker(5 * time.Second)
		go heartbeat(*ticker, c)
	}
}
