# Buz

[![License](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/silverton-io/buz)
![test](https://github.com/silverton-io/buz/actions/workflows/test.yml/badge.svg)
![lint](https://github.com/silverton-io/buz/actions/workflows/lint.yml/badge.svg)
[![ref](https://pkg.go.dev/badge/github.com/silverton-io/buz/pkg.svg)](https://pkg.go.dev/github.com/silverton-io/buz/pkg)

![buz](img/buzz.png)

# What is Buz?

Buz is a system for multi-protocol event collection, validation, annotation, and delivery.

It ships as a single lightweight binary for deployment flexibility.

Toss Buz into an [AWS lambda function](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=snowplow&source=aws-lambda&location=README&rto=https://aws.amazon.com/lambda/), [GCP Cloud Run](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=gcp-cloud-run&source=github&location=README&rto=https://cloud.google.com/run) service, or K8s pod and start collecting events in minutes.

## Multiple Input Protocols

Buz supports [multiple input protocols](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=multiple-event-protocols&source=github&location=README&rto=https://buz.dev/inputs/overview) including:

* [Snowplow Analytics](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=snowplow-analytics&source=github&location=README&rto=https://buz.dev/inputs/saas/snowplow)
* [Cloudevents](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=cloudevents&source=github&location=README&rto=https://buz.dev/inputs/cloudNative/cloudevents)
* [Self-describing JSON](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=self-describing-json&source=github&location=README&rto=https://buz.dev/inputs/buz/self-describing)
* [Webhooks](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=webhooks&source=github&location=README&rto=https://buz.dev/inputs/buz/webhook)

It even hosts a [pixel](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=pixel&source=github&location=README&rto=https://buz.dev/inputs/buz/pixel) for use in constrained tracking environments.

SDK's are supported out of the box so you can point existing Snowplow Analytics or Cloudevents tracking directly to Buz and it will just work™.

Or point your in-house SDK to it using self-describing JSON.

## Multiple (Simultaneous) Destinations

Buz supports [two dozen different destinations](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=destination-overview&source=github&location=README&rto=https://buz.dev/outputs/overview) including:
* [Redpanda](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=dest-redpanda&source=github&location=README&rto=https://buz.dev/outputs/stream/redpanda)
* [Postgres](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=dest-postgres&source=github&location=README&rto=https://buz.dev/outputs/database/postgres)
* [Kinesis Firehose](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=dest-firehose&source=github&location=README&rto=https://buz.dev/outputs/stream/aws-kinesis-firehose)
* [Google Pub/Sub](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=dest-pubsub&source=github&location=README&rto=https://buz.dev/outputs/stream/google-pubsub)
* [Splunk](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=dest-splunk&source=github&location=README&rto=https://buz.dev/outputs/database/splunk)
* [AWS EventBridge](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=dest-eventbridge&source=github&location=README&rto=https://buz.dev/outputs/messageBus/aws-eventbridge)
* [TimescaleDB](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=dest-timescale&source=github&location=README&rto=https://buz.dev/outputs/timeseries/timescaledb)
* [NATS](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=dest-nats&source=github&location=README&rto=https://buz.dev/outputs/messageBus/nats)
* ...and many more.

You can send events to **one or more** destinations, so fanning them out to where they need to Bee is simple. As is using Buz to migrate from one destination system to another.

## Jsonschema-based validation

Every incoming payload is validated in microseconds using [JSON Schema](https://json-schema.org/).

If a payload doesn't conform to the associated schema, it is marked as such.

If a payload doesn't have an explicitly-associated schema (such as the case with arbitrary webhooks and pixels), payload contents are not validated. It is enveloped as `arbitrary` for downstream processing.

## Onboard Schema Registry

Buz ships with a lightweight schema registry that supports [multiple schema backends](https://buz.dev/schema-registry/overview) including:

* [GCS](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=registry-gcs&source=github&location=README&rto=https://buz.dev/schema-registry/backends/object/gcs)
* [S3](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=registry-s3&source=github&location=README&rto=https://buz.dev/schema-registry/backends/object/s3)
* [Minio](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=registry-minio&source=github&location=README&rto=https://buz.dev/schema-registry/backends/object/minio)
* [Postgres](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=registry-postgres&source=github&location=README&rto=https://buz.dev/schema-registry/backends/database/postgres)
* [Mysql](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=registry-mysql&source=github&location=README&rto=https://buz.dev/schema-registry/backends/database/mysql)
* [Mongodb](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=registry-mongodb&source=github&location=README&rto=https://buz.dev/schema-registry/backends/database/mongodb)
* [Local filesystem](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=registry-filesystem&source=github&location=README&rto=https://buz.dev/schema-registry/backends/buz/filesystem)
* ..and more

Schemas are [cached locally](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=registry-cache&source=github&location=README&rto=https://buz.dev/schema-registry/overview#onboard-schema-registry-cache) once sourced from the configured backend. Cache ttl and maximum size are configurable bu have sane defaults.

Schemas are available via HTTP at `/s/$PATH_TO_SCHEMA` or `/s/$SCHEMA_NAME`, depending on the backend.

The onboard schema cache can be purged via a `GET` or `POST` to the `/c/purge` route.

## Payload Enveloping

Each incoming payload is wrapped in a lightweight envelope.

This envelope appends a bit of metadata such as `isValid`, `buzTimestamp`, schema `vendor`, `namespace`, `version`, an event `uuid`, the associated `protocol`, etc. Metadata is then used to power payload **routing** and **sharding** as well as **namespace-level statistics**.

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

# Why Buz?

### It's lightweight


It minimizes the typical infrastructure footprint of collecting events from many different sources and allows for low-infrastructure, highly-flexible rollouts.

### It's flexible


Buz doesn't care what your existing systems look like or what you want them to look like in the future.

It helps with the "now", and helps get your infrastructure to where you'd like it to be (without another migration).

### It saves time and money


Buz aims to **improve the lives of pipeline maintainers** and **drastically reduce long-term maintenance of event collection systems.**


Roll it out fast, keep it going without much thought, and shut it off when it isn't doing anything.

# Try it out

(No, you don't need to talk to anyone. Though we're relatively friendly and there's a [Discord](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=discord&source=github&location=README&rto=https://discord.com/invite/JFKVnVdF2m) if you want to...)

You'll need [go](https://docstream.buz.dev/pixel/io.silverton/buz/pixel/linkClick/v1.0?link=golang&source=github&location=README&rto=https://go.dev/doc/install) on your machine. But don't need to be a [gopher](https://go.dev/blog/gopher).


## Bootstrapping Buz

**Clone:**

    $ git clone git@github.com:silverton-io/buz.git && cd buz


**Bootstrap:**

    $ make bootstrap


## Sending sample events

Events will be sent to two sinks by default - colorized envelopes will be sent to `stdout` and sent to `buz_events.json` or `buz_invalid_events.json` files.

### POST a cloudevent

    curl -X POST localhost:8080/cloudevents -H 'Content-Type:application/cloudevents+json' -d '{"dataschema":"io.silverton/buz/example/gettingStarted/v1.0.json", "data": {"userId": 10, "name": "you", "action": "didSomething"}}'

### POST an arbitrary webhook

    curl -X POST "localhost:8080/webhook" -H 'Content-Type:application/json' -d '{"arbitrary": "thing"}'


### POST a named (schematized) webhook

    curl -X POST "localhost:8080/webhook/io.silverton/buz/example/generic/sample/v1.0" -H 'Content-Type:application/json' -d '{"id": "10"}'


### GET an arbitrary pixel

    curl "localhost:8080/pixel?msg=hello&subject=world"


### GET a named (schematized) pixel

    curl "localhost:8080/pixel/io.silverton/buz/example/generic/sample/v1.0?id=10"


### POST self-describing JSON

    curl -X POST "localhost:8080/self-describing" -H 'Content-Type:application/json' -d '{"payload": {"schema": "io.silverton/buz/example/generic/sample/v1.0.json", "data": {"id": "10"}}}'



# Buz plays nicely with others

Quickstart documentation for setting up a lightweight streaming stack with Buz, a sample ui, nginx, Redpanda, and Kowl can [be found here](https://buz.dev/examples/quickstart).


# Deploying Buz

Buz can be deployed in a [variety of ways](https://buz.dev/deploying/overview). We've included end-to-end (terraformed) samples for AWS and GCP:

* [Buz on AWS Lambda](https://buz.dev/deploying/aws/lambda)
* [Buz on GCP Cloud Run](https://buz.dev/deploying/gcp/cloud_run)

# Documentation

Full documentation can [be found here](https://buz.dev).


🍻🐝
