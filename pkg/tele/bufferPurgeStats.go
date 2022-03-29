package tele

import (
	"sync"
	"time"
)

type BufferPurgeStats struct {
	vmu        sync.Mutex
	imu        sync.Mutex
	Invalid    int        `json:"invalid"`
	Valid      int        `json:"valid"`
	LastPurged *time.Time `json:"lastPurged"`
}

func (s *BufferPurgeStats) incrementValid() {
	s.vmu.Lock()
	defer s.vmu.Unlock()
	s.Valid = s.Valid + 1
}

func (s *BufferPurgeStats) incrementInvalid() {
	s.imu.Lock()
	defer s.imu.Unlock()
	s.Invalid = s.Invalid + 1
}

func (s *BufferPurgeStats) Increment() {
	s.incrementInvalid()
	s.incrementValid()
	n := time.Now().UTC()
	s.LastPurged = &n
}
