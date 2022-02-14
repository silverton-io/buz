package tele

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/http"
	"github.com/silverton-io/gosnowplow/pkg/snowplow"
	"github.com/silverton-io/gosnowplow/pkg/util"
)

const (
	DEFAULT_HOST    string = "http://localhost:8081/gen/p"
	HEARTBEAT_1_0_0 string = "com.silverton.io.tele/heartbeat/jsonschema/1-0-0"
	SNAPSHOT_1_0_0  string = "com.silverton.io.tele/snapshot/jsonschema/1-0-0"
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
		heartbeatPayload := snowplow.SelfDescribingPayload{
			Schema: HEARTBEAT_1_0_0,
			Data:   data,
		}
		http.SendJson(DEFAULT_HOST, heartbeatPayload)
		log.Debug().Msg("heartbeat")
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
		snapshotPayload := snowplow.SelfDescribingPayload{
			Schema: SNAPSHOT_1_0_0,
			Data:   data,
		}
		http.SendJson(DEFAULT_HOST, snapshotPayload)
		ticker := time.NewTicker(5 * time.Second)
		go heartbeat(*ticker, c)
	}
}
