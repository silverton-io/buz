package sink

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/protocol"
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

func (ms *MockSink) BatchPublishValidAndInvalid(ctx context.Context, inputType string, validEvents []event.Envelope, invalidEvents []event.Envelope, meta *tele.Meta) {
	ms.Called(ctx, inputType, validEvents, invalidEvents, meta)
}

func (ms *MockSink) Close() {
	ms.Called()
}

func TestBuildSink(t *testing.T) {
	c := config.Sink{
		Type:                  PUBSUB,
		Project:               "myproject",
		Brokers:               []string{"broker1"},
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
		Brokers:               []string{"broker1"},
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

func TestIncrementStats(t *testing.T) {
	id := uuid.New()
	n := time.Now()
	m := tele.Meta{
		Version:                        "1.x.x",
		InstanceId:                     id,
		StartTime:                      n,
		TrackerDomain:                  "somewhere.something",
		CookieDomain:                   "elsewhere.somewhere.something",
		ValidSnowplowEventsProcessed:   0,
		InvalidSnowplowEventsProcessed: 0,
		ValidGenericEventsProcessed:    0,
		InvalidGenericEventsProcessed:  0,
		ValidCloudEventsProcessed:      0,
		InvalidCloudEventsProcessed:    0,
	}

	incrementStats(protocol.CLOUDEVENTS, 1, 1, &m)
	incrementStats(protocol.SNOWPLOW, 1, 1, &m)
	incrementStats(protocol.GENERIC, 1, 1, &m)
	incrementStats("other", 1, 1, &m)

	assert.Equal(t, int64(1), m.ValidCloudEventsProcessed)
	assert.Equal(t, int64(1), m.InvalidCloudEventsProcessed)
	assert.Equal(t, int64(1), m.ValidGenericEventsProcessed)
	assert.Equal(t, int64(1), m.InvalidGenericEventsProcessed)
	assert.Equal(t, int64(2), m.ValidSnowplowEventsProcessed)
	assert.Equal(t, int64(2), m.InvalidSnowplowEventsProcessed)
}
