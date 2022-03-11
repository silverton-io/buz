package util

import "time"

func GetDuration(start time.Time, end time.Time) time.Duration {
	return end.Sub(start)
}
