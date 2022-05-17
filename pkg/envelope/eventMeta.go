package envelope

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/google/uuid"
)

type EventMeta struct {
	Protocol           string    `json:"protocol"`
	Uuid               uuid.UUID `json:"uuid"`
	Vendor             string    `json:"vendor"`
	PrimaryNamespace   string    `json:"primaryNamespace"`
	SecondaryNamespace string    `json:"secondaryNamespace"`
	TertiaryNamespace  string    `json:"tertiaryNamespace"`
	Name               string    `json:"name"`
	Version            string    `json:"version"`
	Format             *string   `json:"format"`
	Path               string    `json:"path"`
}

func (e EventMeta) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e EventMeta) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}
