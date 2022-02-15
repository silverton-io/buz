package tele

import (
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/http"
	"github.com/silverton-io/gosnowplow/pkg/snowplow"
	"github.com/silverton-io/gosnowplow/pkg/util"
)

const (
	DEFAULT_HOST    string = "http://some.where.else:8081/gen/p"
	HEARTBEAT_1_0_0 string = "com.silverton.io.tele/heartbeat/jsonschema/1-0-0"
	SNAPSHOT_1_0_0  string = "com.silverton.io.tele/snapshot/jsonschema/1-0-0"
)

type Meta struct {
	Version    string    `json:"version"`
	InstanceId uuid.UUID `json:"instanceId"`
	StartTime  time.Time `json:"startTime"`
	Domain     string    `json:"domain"`
}

func (m *Meta) elapsed() float64 {
	return time.Since(m.StartTime).Seconds()
}

type configSnapshot struct {
	GosnowplowVersion string        `json:"gosnowplowVersion"`
	InstanceId        uuid.UUID     `json:"instanceId"`
	Domain            string        `json:"domain"`
	Time              time.Time     `json:"time"`
	Config            config.Config `json:"config"`
}

type beat struct {
	GosnowplowVersion string    `json:"gosnowplowVersion"`
	InstanceId        uuid.UUID `json:"instanceId"`
	Domain            string    `json:"domain"`
	Time              time.Time `json:"time"`
	ElapsedSeconds    float64   `json:"elapsedSeconds"`
}

func heartbeat(t time.Ticker, m *Meta) {
	for _ = range t.C {
		b := beat{
			GosnowplowVersion: m.Version,
			InstanceId:        m.InstanceId,
			Domain:            m.Domain,
			Time:              time.Now(),
			ElapsedSeconds:    m.elapsed(),
		}
		data := util.StructToMap(b)
		heartbeatPayload := snowplow.SelfDescribingPayload{
			Schema: HEARTBEAT_1_0_0,
			Data:   data,
		}
		http.SendJson(DEFAULT_HOST, heartbeatPayload)
	}
}

func Metry(c *config.Config, m *Meta) {
	if c.Tele.Enable {
		snapshot := configSnapshot{
			GosnowplowVersion: m.Version,
			InstanceId:        m.InstanceId,
			Domain:            m.Domain,
			Time:              time.Now(),
			Config:            *c,
		}
		data := util.StructToMap(snapshot)
		snapshotPayload := snowplow.SelfDescribingPayload{
			Schema: SNAPSHOT_1_0_0,
			Data:   data,
		}
		http.SendJson(DEFAULT_HOST, snapshotPayload)
		ticker := time.NewTicker(5 * time.Second)
		go heartbeat(*ticker, m)
	}
}
