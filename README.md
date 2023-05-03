# Buz

[![License](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/silverton-io/buz)

<!-- ![tests](https://github.com/silverton-io/buz/actions/workflows/test/badge.svg) -->

![buz](img/buzz.png)

# What is Buz?

Buz is a system for multi-protocol event collection, validation, annotation, and delivery.

It ships as a single lightweight binary for deployment flexibility.

Toss Buz into an [AWS lambda function](https://aws.amazon.com/lambda/), [GCP Cloud Run](https://cloud.google.com/run) service, or K8s pod and start collecting events in minutes.

## Multiple Input Protocols

Buz supports [multiple input protocols](https://buz.dev/inputs/overview) including:

* [Snowplow Analytics](https://buz.dev/inputs/saas/snowplow)
* [Cloudevents](https://buz.dev/inputs/cloudNative/cloudevents)
* [Self-describing JSON](https://buz.dev/inputs/buz/self-describing)
* [Webhooks](https://buz.dev/inputs/buz/webhook)

It even hosts a [pixel](https://buz.dev/inputs/buz/pixel) for use in constrained tracking environments.

SDK's are supported out of the box so you can point existing Snowplow Analytics or Cloudevents tracking directly to Buz and it will just work™.

Or point your in-house SDK to it using self-describing JSON.

## Multiple (Simultaneous) Destinations

Buz supports [two dozen different destinations](https://buz.dev/outputs/overview) including:
* [Redpanda](https://buz.dev/outputs/stream/redpanda)
* [Postgres](https://buz.dev/outputs/database/postgres)
* [Kinesis Firehose](https://buz.dev/outputs/stream/aws-kinesis-firehose)
* [Google Pub/Sub](https://buz.dev/outputs/stream/google-pubsub)
* [Splunk](https://buz.dev/outputs/database/splunk)
* [AWS EventBridge](https://buz.dev/outputs/messageBus/aws-eventbridge)
* [TimescaleDB](https://buz.dev/outputs/timeseries/timescaledb)
* [NATS](https://buz.dev/outputs/messageBus/nats)
* ...and many more.

You can send events to **one or more** destinations, so fanning them out to where they need to Bee is simple. As is using Buz to migrate from one destination system to another.

## Jsonschema-based validation

Every incoming payload is validated in microseconds using [JSON Schema](https://json-schema.org/).

If a payload doesn't conform to the associated schema, it is marked as such.

If a payload doesn't have an explicitly-associated schema (such as the case with arbitrary webhooks and pixels), payload contents are not validated. It is enveloped as `arbitrary` for downstream processing.

## Onboard Schema Registry

Buz ships with a lightweight schema registry that supports [multiple schema backends](https://buz.dev/schema-registry/overview) including:

* [GCS](https://buz.dev/schema-registry/backends/object/gcs)
* [S3](https://buz.dev/schema-registry/backends/object/s3)
* [Minio](https://buz.dev/schema-registry/backends/object/minio)
* [Postgres](https://buz.dev/schema-registry/backends/database/postgres)
* [Mysql](https://buz.dev/schema-registry/backends/database/mysql)
* [Mongodb](https://buz.dev/schema-registry/backends/database/mongodb)
* [Local filesystem](https://buz.dev/schema-registry/backends/buz/filesystem)
* ..and more

Schemas are [cached locally](https://buz.dev/schema-registry/overview#onboard-schema-registry-cache) once sourced from the configured backend. The cache ttl and maximum size is configurable, but ships with sane defaults.

Schemas are available via HTTP at `/s/$PATH_TO_SCHEMA` or `/s/$SCHEMA_NAME`, depending on the backend.

The onboard schema cache can be purged via a `GET` or `POST` to the `/c/purge` route.

## Payload Enveloping

Each incoming payload is wrapped in a lightweight envelope. This envelope appends metadata such as `isValid`, `buzTimestamp`, schema `vendor`, `namespace`, `version`, an event `uuid`, the associated `protocol`, etc.

Envelope metadata is used to power **routing** and **sharding** far downstream of collection, as well as **namespace-level statistics**.

**As an example of an `arbitrary pixel` event, wrapped in said envelope:**

```
{
    "uuid": "1f9a7a20-8fa7-4179-a0c2-35a80783854a",
    "timestamp": "2023-05-03T02:50:59.464042Z",
    "buzTimestamp": "2023-05-03T02:50:59.464042Z",
    "buzVersion": "x.x.dev",
    "buzName": "buz-bootstrap",
    "buzEnv": "development",
    "protocol": "pixel",
    "schema": "io.silverton/buz/pixel/arbitrary/v1.0.json",
    "vendor": "io.silverton",
    "namespace": "buz.pixel.arbitrary",
    "version": "1.0",
    "isValid": true,
    "payload": {
        "msg": "hello",
        "subject": "zander"
    }
}
```

## Time and Cost Efficiences

Buz aims to **improve the lives of pipeline maintainers** and **drastically reduce long-term maintenance of event collection systems.**

It minimizes the typical infrastructure footprint of collecting events from many different sources and allows for low-infrastructure, highly-flexible rollouts.

When deployed in serverless environments such as AWS Lambda or GCP Cloud Run, Buz is able to follow the utilization curve closely. Which drastically minimizes cost.


# Quickstart

Quickstart documentation for setting up a lightweight streaming stack with Buz, Redpanda, and Kowl can [be found here](https://buz.dev/examples/quickstart).


# Documentation

Documentation can [be found here](https://buz.dev).
