package sink

import (
	"context"
	"errors"
	"testing"

	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/tele"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSink struct {
	mock.Mock
}

func (ms *MockSink) Initialize(conf config.Sink) {
	ms.Called(conf)
}

func (ms *MockSink) BatchPublishValid(ctx context.Context, validEnvelopes []envelope.Envelope) {
	ms.Called()
}

func (ms *MockSink) BatchPublishInvalid(ctx context.Context, invalidEnvelopes []envelope.Envelope) {
	ms.Called()
}

func (ms *MockSink) BatchPublishValidAndInvalid(ctx context.Context, validEvents []envelope.Envelope, invalidEvents []envelope.Envelope, meta *tele.Meta) {
	ms.Called(ctx, validEvents, invalidEvents, meta)
}

func (ms *MockSink) Close() {
	ms.Called()
}

func TestBuildSink(t *testing.T) {
	c := config.Sink{
		Type:                  PUBSUB,
		Project:               "myproject",
		KafkaBrokers:          []string{"broker1"},
		ValidEventTopic:       "valid-topic",
		InvalidEventTopic:     "invalid-topic",
		BufferByteThreshold:   1,
		BufferRecordThreshold: 1,
		BufferDelayThreshold:  1,
	}

	t.Run(PUBSUB, func(t *testing.T) {
		sink, err := BuildSink(c)
		pubsubSink := PubsubSink{}
		assert.IsType(t, &pubsubSink, sink)
		assert.Equal(t, err, nil)
	})

	t.Run(KAFKA, func(t *testing.T) {
		c.Type = KAFKA
		sink, err := BuildSink(c)
		kafkaSink := KafkaSink{}
		assert.IsType(t, &kafkaSink, sink)
		assert.Equal(t, nil, err)
	})

	t.Run(KINESIS, func(t *testing.T) {
		c.Type = KINESIS
		sink, err := BuildSink(c)
		kinesisSink := KinesisSink{}
		assert.IsType(t, &kinesisSink, sink)
		assert.Equal(t, nil, err)
	})

	t.Run(KINESIS_FIREHOSE, func(t *testing.T) {
		c.Type = KINESIS_FIREHOSE
		sink, err := BuildSink(c)
		firehoseSink := KinesisFirehoseSink{}
		assert.IsType(t, &firehoseSink, sink)
		assert.Equal(t, nil, err)
	})

	t.Run(STDOUT, func(t *testing.T) {
		c.Type = STDOUT
		sink, err := BuildSink(c)
		stdoutSink := StdoutSink{}
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
		Type:                  PUBSUB,
		Project:               "myproject",
		KafkaBrokers:          []string{"broker1"},
		ValidEventTopic:       "valid-topic",
		InvalidEventTopic:     "invalid-topic",
		BufferByteThreshold:   1,
		BufferRecordThreshold: 1,
		BufferDelayThreshold:  1,
	}
	mSink := MockSink{}
	mSink.On("Initialize", c)

	InitializeSink(c, &mSink)

	mSink.AssertCalled(t, "Initialize", c)
}
