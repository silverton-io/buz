---
sidebar_position: 1
slug: /
title: What is Honeypot?
---

Honeypot is a system for collecting events from various sources, validating data quality, and delivering them to where they need to bee.

It is designed to be ***easily-configured***, ***easily-deployed***, and ***easily-maintained***. Yet uncompromising with its speed, guarantees, and operational flexibility.


:::tip Quickstart
To rapidly bootstrap a streaming stack using Honeypot, [Redpanda](https://github.com/redpanda-data/), [Kowl](https://github.com/cloudhut/kowl/), and [Materialize](https://github.com/MaterializeInc/materialize) please see the [Quickstart](/examples/quickstart)!
:::


## Consists of a Single Self-Contained Binary

Event collection infrastructure comes in all shapes and sizes except `a single artifact that can be deployed anywhere`. Honeypot changes that.


**Honeypot was created from the ground-up to eliminate as many moving pieces as possible without sacrificing quality and transport guarantees.**


:::tip TL;DR
- Less headaches for infrastructure humans.
- Decreased infrastructure costs
- Increased operational efficiencies.
:::



## Collects Events from One or More Sources

Most event collection systems are single-protocol, which results in duplicate infrastructure when more than one protocol must be collected  -> `webhook`, `snowplow`, `segment`, `cloudevents`, etc.

This is not the case with Honeypot.

**Honeypot supports a number of common input protocols and will continue to support more.**

:::tip TL;DR
- Single, flexible system instead of N pipelines to support N protocols.
:::

## Validates Events on the Edge

Streaming data is all about one thing: increasing the speed of action and decision-making. If events are not validated fast, decisions and actions cannot be made with the conviction they require. Who cares about making bad decisions, on bad data, fast? Nobody.

:::tip TL;DR
- Incoming events are validated immediately.
:::

## Sends Events to One or More Destinations

Honeypot was purpose-built for flexibility and does not require provisioning *more* infrastructure to fan out events to multiple places.

Multiple Honeypot sinks can be configured regardless of what cloud the destination actually resides in.

:::tip TL;DR
- Send events to one place or five. No additional infrastructure necessary.
:::

## Generates Rich Metadata

After years of building, maintaining, and managing event tracking systems there's one thing that has consistently stuck out:

**More metadata and expedited knowledge of what is happening within the stream would be unbelievably empowering**.

Metadata often happens at the end of the pipeline, in the data warehouse. **Not awesome.**

:::tip TL;DR
- More visibility into events, faster.
:::


## Eases Operational Burden

Honeypot is easily configured with a single `yml` file. This file is self-validating and, when using [the vcode yaml plugin](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml), practically writes itself.

It can be rapidly deployed and re-configured, and scales well as needs do.
