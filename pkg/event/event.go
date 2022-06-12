package event

import "database/sql/driver"

type Event interface {
	SchemaName() *string
	PayloadAsByte() ([]byte, error)
	Value() (driver.Value, error)
	Scan(input interface{}) error
}
