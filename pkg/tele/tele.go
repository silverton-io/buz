package tele

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/gosnowplow/pkg/config"
	"github.com/silverton-io/gosnowplow/pkg/event"
	h "github.com/silverton-io/gosnowplow/pkg/http"
	"github.com/silverton-io/gosnowplow/pkg/util"
)

const (
	DEFAULT_ENDPOINT string = "http://some.where.else:8081/gen/p"
	STARTUP_1_0_0    string = "com.silverton.io.tele/startup/1-0-0"
	HEARTBEAT_1_0_0  string = "com.silverton.io.tele/heartbeat/1-0-0"
	SHUTDOWN_1_0_0   string = "com.silverton.io.tele/shutdown/1-0-0"
)

type Meta struct {
	Version                        string    `json:"version"`
	InstanceId                     uuid.UUID `json:"instanceId"`
	StartTime                      time.Time `json:"startTime"`
	TrackerDomain                  string    `json:"trackerDomain"`
	CookieDomain                   string    `json:"cookieDomain"`
	ValidSnowplowEventsProcessed   int64     `json:"validSnowplowEventsProcessed"`
	InvalidSnowplowEventsProcessed int64     `json:"invalidSnowplowEventsProcessed"`
	ValidGenericEventsProcessed    int64     `json:"validGenericEventsProcessed"`
	InvalidGenericEventsProcessed  int64     `json:"invalidGenericEventsProcessed"`
}

func (m *Meta) elapsed() float64 {
	return time.Since(m.StartTime).Seconds()
}

type startup struct {
	Meta   *Meta         `json:"meta"`
	Time   time.Time     `json:"time"`
	Config config.Config `json:"config"`
}

type beat struct {
	Meta           *Meta     `json:"meta"`
	Time           time.Time `json:"time"`
	ElapsedSeconds float64   `json:"elapsedSeconds"`
}

type shutdown struct {
	Meta           *Meta     `json:"meta"`
	Time           time.Time `json:"time"`
	ElapsedSeconds float64   `json:"elapsedSeconds"`
}

func heartbeat(t time.Ticker, m *Meta) {
	for _ = range t.C {
		log.Trace().Msg("sending heartbeat telemetry")
		b := beat{
			Meta:           m,
			Time:           time.Now(),
			ElapsedSeconds: m.elapsed(),
		}
		data := util.StructToMap(b)
		heartbeatPayload := event.SelfDescribingEnvelope{
			Contexts: nil,
			Event: event.SelfDescribingPayload{
				Schema: HEARTBEAT_1_0_0,
				Data:   data,
			},
		}
		endpoint, _ := url.Parse(DEFAULT_ENDPOINT)
		h.SendJson(*endpoint, heartbeatPayload)
	}
}

func Sis(m *Meta) {
	log.Trace().Msg("sending shutdown telemetry")
	shutdown := shutdown{
		Meta:           m,
		Time:           time.Now(),
		ElapsedSeconds: m.elapsed(),
	}
	data := util.StructToMap(shutdown)
	shutdownPayload := event.SelfDescribingEnvelope{
		Contexts: nil,
		Event: event.SelfDescribingPayload{
			Schema: SHUTDOWN_1_0_0,
			Data:   data,
		},
	}
	endpoint, _ := url.Parse(DEFAULT_ENDPOINT)
	h.SendJson(*endpoint, shutdownPayload)
}

func Metry(c *config.Config, m *Meta) {
	if c.Tele.Enabled {
		log.Trace().Msg("sending startup telemetry")
		startup := startup{
			Meta:   m,
			Time:   time.Now(),
			Config: *c,
		}
		data := util.StructToMap(startup)
		startupPayload := event.SelfDescribingEnvelope{
			Contexts: nil,
			Event: event.SelfDescribingPayload{
				Schema: STARTUP_1_0_0,
				Data:   data,
			},
		}
		endpoint, _ := url.Parse(DEFAULT_ENDPOINT)
		h.SendJson(*endpoint, startupPayload)
		ticker := time.NewTicker(5 * time.Second)
		go heartbeat(*ticker, m)
	}
}
