package tele

import (
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/request"
	"github.com/silverton-io/honeypot/pkg/util"
)

const (
	DEFAULT_ENDPOINT string = "https://trck.slvrtnio.com/gen/p"
	STARTUP_1_0      string = "com.silverton.io/honeypot/tele/startup/v1.0.json"
	HEARTBEAT_1_0    string = "com.silverton.io/honeypot/tele/beat/v1.0.json"
	SHUTDOWN_1_0     string = "com.silverton.io/honeypot/tele/shutdown/v1.0.json"
)

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
			Time:           time.Now().UTC(),
			ElapsedSeconds: m.elapsed(),
		}
		data := util.StructToMap(b)
		heartbeatPayload := event.SelfDescribingEvent{
			Contexts: nil,
			Payload: event.SelfDescribingPayload{
				Schema: HEARTBEAT_1_0,
				Data:   data,
			},
		}
		endpoint, _ := url.Parse(DEFAULT_ENDPOINT)
		request.PostEvent(*endpoint, heartbeatPayload)
	}
}

func Sis(m *Meta) {
	log.Trace().Msg("sending shutdown telemetry")
	shutdown := shutdown{
		Meta:           m,
		Time:           time.Now().UTC(),
		ElapsedSeconds: m.elapsed(),
	}
	data := util.StructToMap(shutdown)
	shutdownPayload := event.SelfDescribingEvent{
		Contexts: nil,
		Payload: event.SelfDescribingPayload{
			Schema: SHUTDOWN_1_0,
			Data:   data,
		},
	}
	endpoint, _ := url.Parse(DEFAULT_ENDPOINT)
	request.PostEvent(*endpoint, shutdownPayload)
}

func Metry(c *config.Config, m *Meta) {
	if c.Tele.Enabled {
		log.Trace().Msg("sending startup telemetry")
		startup := startup{
			Meta:   m,
			Time:   time.Now().UTC(),
			Config: *c,
		}
		data := util.StructToMap(startup)
		startupPayload := event.SelfDescribingEvent{
			Contexts: nil,
			Payload: event.SelfDescribingPayload{
				Schema: STARTUP_1_0,
				Data:   data,
			},
		}
		endpoint, _ := url.Parse(DEFAULT_ENDPOINT)
		request.PostEvent(*endpoint, startupPayload)
		ticker := time.NewTicker(time.Duration(30000) * time.Millisecond)
		go heartbeat(*ticker, m)
	}
}
