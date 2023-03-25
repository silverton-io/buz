// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package stats

import (
	"sync"
)

type ProtocolStats struct {
	vmu     sync.Mutex
	imu     sync.Mutex
	Invalid map[string]map[string]int64 `json:"invalid"`
	Valid   map[string]map[string]int64 `json:"valid"`
}

// func (ps *ProtocolStats) Build() {
// 	var vProtoStat = make(map[string]map[string]int64)
// 	var invProtoStat = make(map[string]map[string]int64)
// 	ps.Valid = vProtoStat
// 	ps.Invalid = invProtoStat
// 	for _, protocol := range protocol.GetInputProtocols() {
// 		var vEventStat = make(map[string]int64)
// 		var invEventStat = make(map[string]int64)
// 		ps.Valid[protocol] = vEventStat
// 		ps.Invalid[protocol] = invEventStat
// 	}
// }

// func (ps *ProtocolStats) IncrementValid(event *envelope.EventMeta, count int64) {
// 	i := ps.Valid[event.Protocol][event.Namespace]
// 	ps.vmu.Lock()
// 	defer ps.vmu.Unlock()
// 	ps.Valid[event.Protocol][event.Namespace] = i + count
// }

// func (ps *ProtocolStats) IncrementInvalid(event *envelope.EventMeta, count int64) {
// 	i := ps.Invalid[event.Protocol][event.Namespace]
// 	ps.imu.Lock()
// 	defer ps.imu.Unlock()
// 	ps.Invalid[event.Protocol][event.Namespace] = i + count
// }

// func BuildProtocolStats() *ProtocolStats {
// 	ps := ProtocolStats{}
// 	ps.Build()
// 	return &ps
// }
