// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package util

import "time"

func GetDuration(start time.Time, end time.Time) time.Duration {
	return end.Sub(start)
}
