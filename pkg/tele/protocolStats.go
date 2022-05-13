package tele

import (
	"sync"

	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/protocol"
)

type ProtocolStats struct {
	vmu     sync.Mutex
	imu     sync.Mutex
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

func (ps *ProtocolStats) IncrementValid(event *envelope.Event, count int64) {
	i := ps.Valid[event.Protocol][event.Path]
	ps.vmu.Lock()
	defer ps.vmu.Unlock()
	ps.Valid[event.Protocol][event.Path] = i + count
}

func (ps *ProtocolStats) IncrementInvalid(event *envelope.Event, count int64) {
	i := ps.Invalid[event.Protocol][event.Path]
	ps.imu.Lock()
	defer ps.imu.Unlock()
	ps.Invalid[event.Protocol][event.Path] = i + count
}

func BuildProtocolStats() *ProtocolStats {
	ps := ProtocolStats{}
	ps.Build()
	return &ps
}
