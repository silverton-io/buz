package db

import (
	"time"
)

type BasePKeylessModel struct {
	CreatedAt time.Time  `json:"-" sql:"index"`
	UpdatedAt time.Time  `json:"-" sql:"index"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
