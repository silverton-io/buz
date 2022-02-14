package tele

import (
	"fmt"
	"time"

	"github.com/silverton-io/gosnowplow/pkg/config"
)

const (
	DEFAULT_HOST string = "someendpoint"
)

type configSnapshot struct {
	version string        `json:"version"`
	domain  string        `json:"domain"`
	config  config.Config `json:"config"`
}

type beat struct {
	version string    `json:"version"`
	domain  string    `json:"domain"`
	time    time.Time `json:"time"`
}

func heartbeat(t time.Ticker, c config.Config) {
	for t := range t.C {
		fmt.Printf("tick! %v\n", t)
	}
}

func Metry(c config.Config) {
	if c.Tele.Enable {
		ticker := time.NewTicker(5 * time.Second)
		go heartbeat(*ticker, c)
	}
}
