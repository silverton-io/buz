package tele

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/event"
	"github.com/silverton-io/gosnowplow/pkg/http"
	"github.com/silverton-io/gosnowplow/pkg/util"
)

const (
	DEFAULT_HOST    string = "http://some.where.else:8081/gen/p"
	STARTUP_1_0_0   string = "com.silverton.io.tele/startup/1-0-0"
	HEARTBEAT_1_0_0 string = "com.silverton.io.tele/heartbeat/1-0-0"
	SHUTDOWN_1_0_0  string = "com.silverton.io.tele/shutdown/1-0-0"
)

type Meta struct {
	Version                        string    `json:"version"`
	InstanceId                     uuid.UUID `json:"instanceId"`
	StartTime                      time.Time `json:"startTime"`
	Domain                         string    `json:"domain"`
	ValidSnowplowEventsProcessed   int64     `json:"validSnowplowEventsProcessed"`
	InvalidSnowplowEventsProcessed int64     `json:"invalidSnowplowEventsProcessed"`
	ValidGenericEventsProcessed    int64     `json:"validGenericEventsProcessed"`
	InvalidGenericEventsProcessed  int64     `json:"invalidGenericEventsProcessed"`
}

func (m *Meta) elapsed() float64 {
	return time.Since(m.StartTime).Seconds()
}

type startup struct {
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

type shutdown struct {
	GosnowplowVersion string    `json:"gosnowplowVersion"`
	InstanceId        uuid.UUID `json:"instanceId"`
	Domain            string    `json:"domain"`
	Time              time.Time `json:"time"`
	ElapsedSeconds    float64   `json:"elapsedSeconds"`
}

func heartbeat(t time.Ticker, m *Meta) {
	log.Trace().Msg("sending heartbeat")
	for _ = range t.C {
		b := beat{
			GosnowplowVersion: m.Version,
			InstanceId:        m.InstanceId,
			Domain:            m.Domain,
			Time:              time.Now(),
			ElapsedSeconds:    m.elapsed(),
		}
		data := util.StructToMap(b)
		heartbeatPayload := event.SelfDescribingPayload{
			Schema: HEARTBEAT_1_0_0,
			Data:   data,
		}
		http.SendJson(DEFAULT_HOST, heartbeatPayload)
	}
}

func Sis(m *Meta) {
	log.Trace().Msg("sending shutdown")
	shutdown := shutdown{
		GosnowplowVersion: m.Version,
		InstanceId:        m.InstanceId,
		Domain:            m.Domain,
		Time:              time.Now(),
		ElapsedSeconds:    m.elapsed(),
	}
	data := util.StructToMap(shutdown)
	shutdownPayload := event.SelfDescribingPayload{
		Schema: SHUTDOWN_1_0_0,
		Data:   data,
	}
	http.SendJson(DEFAULT_HOST, shutdownPayload)
}

func Metry(c *config.Config, m *Meta) {
	if c.Tele.Enable {
		log.Trace().Msg("sending startup")
		startup := startup{
			GosnowplowVersion: m.Version,
			InstanceId:        m.InstanceId,
			Domain:            m.Domain,
			Time:              time.Now(),
			Config:            *c,
		}
		data := util.StructToMap(startup)
		startupPayload := event.SelfDescribingPayload{
			Schema: STARTUP_1_0_0,
			Data:   data,
		}
		http.SendJson(DEFAULT_HOST, startupPayload)
		ticker := time.NewTicker(5 * time.Second)
		go heartbeat(*ticker, m)
	}
}
