# Buz

[![License](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/silverton-io/buz)

<!-- ![tests](https://github.com/silverton-io/buz/actions/workflows/test/badge.svg) -->

![buz](img/buzz.png)

# What is Buz?

Buz is a system for multi-protocol event collection, validation, annotation, and delivery.

It ships as a single lightweight binary for deployment flexibility. Toss Buz into an [AWS lambda function](https://aws.amazon.com/lambda/), [GCP Cloud Run](https://cloud.google.com/run) service, or K8s pod and immediately start collecting events.

## Multiple Input Protocols

Buz supports [multiple input protocols](https://buz.dev/inputs/overview) including [Snowplow Analytics](https://buz.dev/inputs/saas/snowplow), [Cloudevents](https://buz.dev/inputs/cloudNative/cloudevents), [Self-describing JSON](https://buz.dev/inputs/buz/self-describing), and [Webhooks](https://buz.dev/inputs/buz/webhook). It even hosts a [Pixel](https://buz.dev/inputs/buz/pixel) for use in constrained tracking environments.

SDK's are supported out of the box so you can point Snowplow Analytics and Cloudevents tracking directly to Buz and it will just work™. Or point your in-house SDK to it using self-describing JSON.

## Multiple (Simultaneous) Destinations

Buz supports [two dozen different destinations](https://buz.dev/outputs/overview) including [Redpanda](https://buz.dev/outputs/stream/redpanda), [Postgres](https://buz.dev/outputs/database/postgres), [Kinesis Firehose](https://buz.dev/outputs/stream/aws-kinesis-firehose), [Google Pub/Sub](https://buz.dev/outputs/stream/google-pubsub), [Splunk](https://buz.dev/outputs/database/splunk), [AWS EventBridge](https://buz.dev/outputs/messageBus/aws-eventbridge), [TimescaleDB](https://buz.dev/outputs/timeseries/timescaledb), and many more.

You can send events to **one or more** destinations.

## Jsonschema-based validation

Every incoming payload is validated in microseconds using [JSON Schema](https://json-schema.org/).

If a payload doesn't conform to the associated schema, it is marked as such.

If a payload doesn't have an associated schema (such as the case with arbitrary webhooks and pixels) payload contents are not validated. But are still enveloped as `arbitrary` for downstream processing.

## Payload Enveloping

Each incoming payload is wrapped in a lightweight envelope. This envelope appends metadata such as `isValid`, `buzTimestamp`, schema `vendor`, `namespace`, `version`, an event `uuid`, the associated `protocol`, etc.

Envelope metadata is used to power routing and sharding far downstream of collection, as well as rich internal statistics.

## Time and Cost Efficiences

Buz aims to **improve the lives of pipeline maintainers** and **drastically reduce long-term maintenance of event collection systems.**

It minimizes the typical infrastructure footprint of collecting events from many different sources and allows for low-infrastructure, highly-flexible rollouts.

When deployed in serverless environments such as AWS Lambda or GCP Cloud Run, Buz is able to follow the utilization curve closely. Which drastically minimizes cost.


# Quickstart

Quickstart documentation for setting up a lightweight streaming stack with Buz, Redpanda, and Kowl can [be found here](https://buz.dev/examples/quickstart).


# Documentation

Documentation can [be found here](https://buz.dev).
