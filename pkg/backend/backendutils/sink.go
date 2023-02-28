// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package backendutils

import (
	"github.com/google/uuid"
)

type SinkMetadata struct {
	Id               *uuid.UUID `json:"id"`
	Name             string     `json:"name"`
	SinkType         string     `json:"sinkType"`
	DeliveryRequired bool       `json:"deliveryRequired"`
}
