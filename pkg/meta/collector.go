package meta

import (
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/config"
)

type CollectorMeta struct {
	Version       string    `json:"version"`
	Name          string    `json:"name"`
	InstanceId    uuid.UUID `json:"instanceId"`
	StartTime     time.Time `json:"startTime"`
	TrackerDomain string    `json:"trackerDomain"`
	CookieDomain  string    `json:"cookieDomain"`
}

func (m *CollectorMeta) Elapsed() float64 {
	return time.Since(m.StartTime).Seconds()
}

func BuildCollectorMeta(version string, conf *config.Config) *CollectorMeta {
	instanceId := uuid.New()
	m := CollectorMeta{
		Version:       version,
		Name:          conf.App.Name,
		InstanceId:    instanceId,
		StartTime:     time.Now().UTC(),
		TrackerDomain: conf.App.TrackerDomain,
		CookieDomain:  conf.Cookie.Domain,
	}
	return &m
}
