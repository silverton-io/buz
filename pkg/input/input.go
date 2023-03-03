// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE
package input

import (
	"github.com/gin-gonic/gin"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/manifold"
	"github.com/silverton-io/buz/pkg/meta"
)

type Input interface {
	Initialize(engine *gin.Engine, manifold *manifold.Manifold, conf *config.Config, metadata *meta.CollectorMeta) error
	Handler(m manifold.Manifold, conf config.Config, metadata *meta.CollectorMeta) gin.HandlerFunc
	SquawkboxHandler(m manifold.Manifold, conf config.Config, metadata *meta.CollectorMeta) gin.HandlerFunc
	EnvelopeBuilder(c *gin.Context, conf *config.Config, metadata *meta.CollectorMeta) []envelope.Envelope
	// Routes() []string
	// Auth() interface{}
}
