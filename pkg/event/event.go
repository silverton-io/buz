package event

import "database/sql/driver"

type Event interface {
	SchemaName() *string
	Protocol() string
	PayloadAsByte() ([]byte, error)
	AsByte() ([]byte, error)
	AsMap() (map[string]interface{}, error)
	Value() (driver.Value, error)
	Scan(input interface{}) error
}
