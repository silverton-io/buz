package db

import (
	"time"
)

type BasePKeylessModel struct {
	CreatedAt time.Time  `sql:"index"`
	UpdatedAt time.Time  `sql:"index"`
	DeletedAt *time.Time `sql:"index"`
}
