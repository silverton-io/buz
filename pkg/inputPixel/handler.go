// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package inputpixel

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/params"
	"github.com/silverton-io/buz/pkg/response"
)

const PX string = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8Xw8AAoMBgDTD2qgAAAAASUVORK5CYII="

func Handler(h params.Handler, m manifold.Manifold) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		envelopes := BuildEnvelopesFromRequest(c, h.Config, h.CollectorMeta)
		err := m.Distribute(envelopes)
		if err != nil {
			c.Header("Retry-After", response.RETRY_AFTER_60)
			c.JSON(http.StatusServiceUnavailable, response.ManifoldDistributionError)
		} else {
			b, _ := base64.StdEncoding.DecodeString(PX)
			c.Data(http.StatusOK, "image/png", b)
		}
	}
	return gin.HandlerFunc(fn)
}
