package event

type Event interface {
	Schema() *string
	Protocol() string
	PayloadAsByte() ([]byte, error)
	AsByte() ([]byte, error)
	AsMap() (map[string]interface{}, error)
}
