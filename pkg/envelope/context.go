// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/util"
)

const HTTP_HEADERS_CONTEXT string = "io.silverton/buz/internal/contexts/httpHeaders/v1.0.json"

type Contexts map[string]interface{}

func (c Contexts) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

func (c Contexts) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &c)
}

func BuildContextsFromRequest(c *gin.Context) Contexts {
	headers := util.HttpHeadersToMap(c)
	context := map[string]interface{}{
		HTTP_HEADERS_CONTEXT: headers,
	}
	return context
}

func (c *Contexts) AsByte() ([]byte, error) {
	return json.Marshal(c)
}
