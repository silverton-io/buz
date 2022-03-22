package tele

import (
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/config"
)

type Meta struct {
	Version          string            `json:"version"`
	InstanceId       uuid.UUID         `json:"instanceId"`
	StartTime        time.Time         `json:"startTime"`
	TrackerDomain    string            `json:"trackerDomain"`
	CookieDomain     string            `json:"cookieDomain"`
	ProtocolStats    *ProtocolStats    `json:"protocolStats"`
	BufferPurgeStats *BufferPurgeStats `json:"bufferPurgeStats"`
}

func (m *Meta) elapsed() float64 {
	return time.Since(m.StartTime).Seconds()
}

func BuildMeta(version string, conf *config.Config) *Meta {
	instanceId := uuid.New()
	ps := ProtocolStats{}
	ps.Build()
	bs := BufferPurgeStats{}
	m := Meta{
		Version:          version,
		InstanceId:       instanceId,
		StartTime:        time.Now(),
		TrackerDomain:    conf.App.TrackerDomain,
		CookieDomain:     conf.Cookie.Domain,
		ProtocolStats:    &ps,
		BufferPurgeStats: &bs,
	}
	return &m
}
