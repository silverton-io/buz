package envelope

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

type Event struct {
	Protocol           string    `json:"protocol"`
	Uuid               uuid.UUID `json:"uuid"`
	Vendor             string    `json:"vendor,omitempty"`
	PrimaryNamespace   string    `json:"primaryNamespace,omitempty"`
	SecondaryNamespace string    `json:"secondaryNamespace,omitempty"`
	TertiaryNamespace  string    `json:"tertiaryNamespace,omitempty"`
	Name               string    `json:"name,omitempty"`
	Version            string    `json:"version,omitempty"`
	Format             string    `json:"format,omitempty"`
	Path               string    `json:"path,omitempty"`
}

func (e Event) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e Event) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}
