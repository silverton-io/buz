// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/stretchr/testify/mock"
)

type MockSink struct {
	mock.Mock
}

func (ms *MockSink) Metadata() backendutils.SinkMetadata {
	id := uuid.New()
	sinkName := "test"
	return backendutils.SinkMetadata{
		Id:               id,
		Name:             sinkName,
		DeliveryRequired: false,
		DefaultOutput:    "somewhere",
		DeadletterOutput: "else",
	}
}

func (ms *MockSink) Initialize(conf config.Sink) error {
	ms.Called(conf)
	return nil
}

func (ms *MockSink) Shutdown() error {
	ms.Called()
	return nil
}
