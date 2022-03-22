package tele

import (
	"sync"
	"sync/atomic"

	"github.com/silverton-io/honeypot/pkg/protocol"
)

type ProtocolStats struct {
	mu      sync.Mutex
	Invalid map[string]map[string]int64 `json:"invalid"`
	Valid   map[string]map[string]int64 `json:"valid"`
}

func (ps *ProtocolStats) Build() {
	var vProtoStat = make(map[string]map[string]int64)
	var invProtoStat = make(map[string]map[string]int64)
	ps.Valid = vProtoStat
	ps.Invalid = invProtoStat
	for _, protocol := range protocol.GetIntputProtocols() {
		var vEventStat = make(map[string]int64)
		var invEventStat = make(map[string]int64)
		ps.Valid[protocol] = vEventStat
		ps.Invalid[protocol] = invEventStat
	}
}

func (ps *ProtocolStats) IncrementValid(protocol string, event string, count int64) {
	i := ps.Valid[protocol][event]
	if i == 0 {
		ps.Valid[protocol][event] = count
	} else {
		atomic.AddInt64(&i, count)
	}
}

func (ps *ProtocolStats) IncrementInvalid(protocol string, event string, count int64) {
	i := ps.Invalid[protocol][event]
	if i == 0 {
		ps.Invalid[protocol][event] = count
	} else {
		atomic.AddInt64(&i, count)
	}
}

func (ps *ProtocolStats) Merge(stats *ProtocolStats) {
	for protocol, pStat := range stats.Valid {
		for event, eStat := range pStat {
			ps.mu.Lock()
			defer ps.mu.Unlock()
			i := ps.Valid[protocol][event]
			ps.Valid[protocol][event] = i + eStat
		}
	}
	for protocol, pStat := range stats.Invalid {
		for event, eStat := range pStat {
			ps.mu.Lock()
			defer ps.mu.Unlock()
			i := ps.Invalid[protocol][event]
			ps.Invalid[protocol][event] = i + eStat
		}
	}
}

func BuildProtocolStats() *ProtocolStats {
	ps := ProtocolStats{}
	ps.Build()
	return &ps
}
