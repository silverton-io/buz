package util

import "time"

func GetDuration(start time.Time) time.Duration {
	return time.Now().Sub(start)
}
