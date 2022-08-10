// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the GPLv3 license, a copy of
// which may be found at https://github.com/silverton-io/honeypot/blob/main/LICENSE

package db

import (
	"time"
)

type BasePKeylessModel struct {
	CreatedAt time.Time  `json:"-" sql:"index"`
	UpdatedAt time.Time  `json:"-" sql:"index"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
