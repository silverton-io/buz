// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

type EventMeta struct {
	Protocol          string    `json:"protocol,omitempty"`
	Uuid              uuid.UUID `json:"uuid,omitempty"`
	Vendor            string    `json:"vendor,omitempty"`
	Namespace         string    `json:"namespace,omitempty"`
	Version           string    `json:"version,omitempty"`
	Format            string    `json:"format,omitempty"`
	Schema            string    `json:"schema,omitempty"`
	DisableValidation bool      `json:"disableValidation,omitempty"`
}

func (e EventMeta) Value() (driver.Value, error) {
	b, err := json.Marshal(e)
	return string(b), err
}

func (e EventMeta) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &e)
}

// Return the desired database schema name using event metadata.
// An event with `com.something` vendor should be directed into
// a `com_something` schema.
func (e *EventMeta) DbSchemaName() string {
	return strings.Replace(e.Vendor, ".", "_", -1)
}

// Return the desired database table name using event metadata.
// An event with a namespace of `some.namespace.something` and
// a version of `1.0` and/or `1.1` should be directed into
// a `some_namespace_something_1` table.
func (e *EventMeta) DbTableName() string {
	ns := strings.Replace(e.Namespace, ".", "_", -1)
	// Only worry about the major version of the schema since minors are backwards-compat
	version := strings.Replace(e.Version, ".", "_", -1)
	majorVersion := strings.Split(version, "_")[0]
	return ns + "_" + majorVersion
}
