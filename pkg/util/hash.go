// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the AGPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(s string) string {
	m := md5.Sum([]byte(s))
	return hex.EncodeToString(m[:])
}
