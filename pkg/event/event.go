package event

import "database/sql/driver"

type Event interface {
	Schema() *string
	Protocol() string
	PayloadAsByte() ([]byte, error)
	AsByte() ([]byte, error)
	AsMap() (map[string]interface{}, error)
	Value() (driver.Value, error)
	Scan(input interface{}) error
}
