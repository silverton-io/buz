package tele

import (
	"sync"
	"time"
)

type BufferPurgeStats struct {
	vmu               sync.Mutex
	imu               sync.Mutex
	Invalid           int        `json:"invalid"`
	Valid             int        `json:"valid"`
	InvalidLastPurged *time.Time `json:"invalidLastPurged"`
	ValidLastPurged   *time.Time `json:"validLastPurged"`
}

func (s *BufferPurgeStats) IncrementValid() {
	s.vmu.Lock()
	defer s.vmu.Unlock()
	s.Valid = s.Valid + 1
	n := time.Now().UTC()
	s.ValidLastPurged = &n
}

func (s *BufferPurgeStats) IncrementInvalid() {
	s.imu.Lock()
	defer s.imu.Unlock()
	s.Invalid = s.Invalid + 1
	n := time.Now().UTC()
	s.InvalidLastPurged = &n
}
