
# Honeypot

![Honeypot](img/honeypot.png)

A lightweight, snowplow-compatible streaming event collection system.

Honeypot is primarily built for flexibility, scalability, and speed.

Secondarily for configuration, deployment, and management ease.


## Supported Event Payloads

- Snowplow analytics ✅
- Cloudevents ✅
- Custom self-describing events ✅


## Supported Schema Cache Backends

- S3 ✅
- GCS ✅
- Filesystem ✅
- Remote HTTP/S ✅
- Kafka Schema Registry ❌


## Supported Sinks

- Kafka/ Redpanda ✅
- Pubsub ✅
- Kinesis ✅
- Kinesis Firehose ✅
- File ✅
- HTTP/S ❌
- Clickhouse ❌
- Postgres ❌
- Firebolt ❌
- PubNub ❌


## Supported Deployment Methods

- K8S ✅
- Knative ✅
- Serverless ✅
- Regular ol' vm's ✅
- Anything else that runs wee little persistent docker containers or go binaries ✅


## Supported Endpoint Configuration

- Snowplow
    - Default Snowplow routes ✅
    - Custom Snowplow routes ✅
    - Configurable/disableable open redirects ✅
- Cloudevents
    - Single POST (`application/cloudevents`) ✅
    - Batch POST (`application/cloudevents-batch`) ✅
- Generic Self-Describing
    - Single POST ✅
    - Batch POST ✅
    - Configurable payload, contexts, and schema keys ✅
- Health
    - Configurable healthcheck route ✅
- Stats
    - Configurable/disableable event stats route ✅
