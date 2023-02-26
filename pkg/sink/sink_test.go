// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package sink

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/silverton-io/buz/pkg/backend/backendutils"
	"github.com/silverton-io/buz/pkg/backend/kafka"
	"github.com/silverton-io/buz/pkg/backend/kinesis"
	"github.com/silverton-io/buz/pkg/backend/kinesisFirehose"
	"github.com/silverton-io/buz/pkg/backend/pubsub"
	"github.com/silverton-io/buz/pkg/backend/stdout"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/constants"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSink struct {
	mock.Mock
}

func (ms *MockSink) Metadata() backendutils.SinkMetadata {
	id := uuid.New()
	sinkName := "test"
	return backendutils.SinkMetadata{
		Id:   &id,
		Name: sinkName,
	}
}

func (ms *MockSink) Initialize(conf config.Sink) error {
	ms.Called(conf)
	return nil
}

func (ms *MockSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) error {
	ms.Called()
	return nil
}

func (ms *MockSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) error {
	ms.Called()
	return nil
}

func (ms *MockSink) BatchPublish(ctx context.Context, envelopes []envelope.Envelope) error {
	ms.Called()
	return nil
}

func (ms *MockSink) BatchPublishValidAndInvalid(ctx context.Context, validEvents []envelope.Envelope, invalidEvents []envelope.Envelope) {
	ms.Called(ctx, validEvents, invalidEvents)
}

func (ms *MockSink) Close() {
	ms.Called()
}

func TestBuildSink(t *testing.T) {
	c := config.Sink{
		Type:         constants.PUBSUB,
		Project:      "myproject",
		KafkaBrokers: []string{"broker1"},
	}

	t.Run(constants.PUBSUB, func(t *testing.T) {
		sink, err := BuildSink(c)
		pubsubSink := pubsub.Sink{}
		assert.IsType(t, &pubsubSink, sink)
		assert.Equal(t, err, nil)
	})

	t.Run(constants.KAFKA, func(t *testing.T) {
		c.Type = constants.KAFKA
		sink, err := BuildSink(c)
		kafkaSink := kafka.Sink{}
		assert.IsType(t, &kafkaSink, sink)
		assert.Equal(t, nil, err)
	})

	t.Run(constants.KINESIS, func(t *testing.T) {
		c.Type = constants.KINESIS
		sink, err := BuildSink(c)
		kinesisSink := kinesis.Sink{}
		assert.IsType(t, &kinesisSink, sink)
		assert.Equal(t, nil, err)
	})

	t.Run(constants.KINESIS_FIREHOSE, func(t *testing.T) {
		c.Type = constants.KINESIS_FIREHOSE
		sink, err := BuildSink(c)
		firehoseSink := kinesisFirehose.Sink{}
		assert.IsType(t, &firehoseSink, sink)
		assert.Equal(t, nil, err)
	})

	t.Run(constants.STDOUT, func(t *testing.T) {
		c.Type = constants.STDOUT
		sink, err := BuildSink(c)
		stdoutSink := stdout.Sink{}
		assert.IsType(t, &stdoutSink, sink)
		assert.Equal(t, nil, err)
	})

	t.Run("unsupported", func(t *testing.T) {
		c.Type = "unsupported-type"
		wantedErr := errors.New("unsupported sink: " + c.Type)
		sink, err := BuildSink(c)
		assert.Equal(t, nil, sink)
		assert.Equal(t, wantedErr, err)
	})
}

func TestInitializeSink(t *testing.T) {
	c := config.Sink{
		Type:         constants.PUBSUB,
		Project:      "myproject",
		KafkaBrokers: []string{"broker1"},
	}
	mSink := MockSink{}
	mSink.On("Initialize", c)

	// InitializeSink(c, &mSink)

	// mSink.AssertCalled(t, "Initialize", c)
}
