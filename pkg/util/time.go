// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import "time"

// GetDuration returns the associated duration between two times
func GetDuration(start time.Time, end time.Time) time.Duration {
	return end.Sub(start)
}
