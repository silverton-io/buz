---
sidebar_position: 1
---

# What is Honeypot?


Honeypot is a multi-protocol event collection, validation, and routing system.

It is designed to be ***easily-configured***, ***easily-deployed***, and ***easily-maintained***.

Yet uncompromising with its speed, guarantees, and operational flexibility. 


:::tip Quickstart
Please see the [Quickstart](/docs/examples/quickstart) to dive head-first into an example of running Honeypot alongside a three-node [Redpanda](https://github.com/redpanda-data/) cluster, [Kowl](https://github.com/cloudhut/kowl/), and [Materialize](https://github.com/MaterializeInc/materialize).
:::


## Consists of a Single Self-Contained Binary

Event collection infrastructure comes in all shapes and sizes except `a single artifact that can be deployed anywhere`.


**Honeypot was created from the ground-up to eliminate as many moving pieces as possible without sacrificing quality and transport guarantees.**


:::tip TL;DR
- Less headaches for infrastructure humans.
- Decreased infrastructure costs
- Increased operational efficiencies.
:::



## Supports Multiple Input Protocols

Single-protocol event collection systems result in duplicate infrastructure that does effectively the same thing -> `webhook`, `snowplow`, `segment`, `cloudevents`, etc.

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


## Is Easily Configured

Honeypot is easily configured with a single `yml` file. This file is self-validating and, when using [the vcode yaml plugin](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml), practically writes itself.

:::tip TL;DR
- You don't need a SaaS product to help you configure your system.
:::
