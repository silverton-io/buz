
# Honeypot

![Honeypot](src/honeypot.png)

A lightweight, snowplow-compatible streaming event collection system.

Honeypot is built for flexibility, scalability, and speed. While simultaneously being easy to configure, deploy, and manage.



## Supported Event Payloads

    - Snowplow analytics ✅
    - Cloudevents ✅
    - Custom self-describing events ✅


## Suppported Schema Cache Backends

    - S3 ✅
    - GCS ✅
    - Filesystem ✅
    - Remote HTTP/S ✅
    - Kafka Schema Registry ✅


## Supported Outputs

    - Kafka/ Redpanda ✅
    - Pubsub ✅
    - Kinesis ❌ (in progress)
    - File ✅


## Supported Deployment Methods

    - K8S ✅
    - Knative ✅
    - Serverless ✅
    - Regular ol' vm's ✅
    - Anything else that runs wee little docker containers ✅


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
