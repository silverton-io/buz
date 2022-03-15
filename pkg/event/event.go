package event

type Event interface {
	Schema() *string
	Protocol() string
	AsByte() ([]byte, error)
	AsMap() (map[string]interface{}, error)
}
