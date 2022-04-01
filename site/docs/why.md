# Why Honeypot?

## **Simplified collection and transport architecture**

Event collection infrastructure comes in all shapes and sizes. These systems are often burdensome to maintain and operate due to dozens and dozens of moving pieces. Or the lack of *important* moving pieces that introduce potential data loss.

**Honeypot was created from the ground-up to eliminate as many moving pieces as possible without sacrificing quality and transport guarantees.**

**TL;DR:**

- Less infrastructure = reduced maintenance and headaches for humans.
- Less infrastructure = reduced infrastructure costs, increased operational efficiencies.


## **Multi-protocol**

Event collection systems are often single-protocol -> think separate systems for collecting/pipelining arbitrary webhooks, Snowplow Analytics, Segment, Cloudevents, etc.

**Honeypot supports a number of common input protocols and will continue to support more.**

**TL;DR:**

- A single, flexible system instead of N pipelines to support N protocols.
- Efficiences of scale.

## **Event Validation on the Edge**

Streaming data is all about one thing: increasing the speed of action and decision-making. If events are not validated fast, decisions and actions cannot be made with the conviction they require. Who cares about making bad decisions, on bad data, fast? Nobody.

**TL;DR:**

- Honeypot validates incoming events the millisecond they are collected.
- Honeypot empowers faster decision-making and action.

## **Multi-cloud**

Honeypot ships with both Go binaries and a Docker image. This results in minimal artifacts to deploy, maintain, and monitor, and an entirely cloud-agnostic deployment model.

Want to deploy on GCP Cloud Run? Easy.

Want to deploy with GKE, EKS, or in k8s? Easy.

Want to deploy on bare metal? Easy.

**TL;DR:**

- Want to have x-cloud event tracking? Honeypot.
- Want to migrate clouds? Honeypot.

## **Rich metadata and stream introspection**

After years of building, maintaining, and managing event tracking systems there's one thing that has always stuck out: ***more metadata and expedited knowledge of what is happening within the stream would be unbelievably empowering***.

This usually happens at the tail end of the pipeline - in the data warehouse. **Honeypot aims to flip this on its head.**

**TL;DR:**

- Honeypot gives more visibility into events, faster.

## **Flexible destinations**

Honeypot was purpose-built to give the practitioner ultimate flexibility. Most event collection systems have a single "destination", or require setting up *more* infrastructure to fan out events to multiple places. We see this as wasteful.

Multiple Honeypot sinks can be configured regardless of what cloud the destination actually resides in.

**TL;DR:**

- Honeypot allows the operator to send events to one place. Or five. No additional infrastructure necessary.


## **Easy configuration**

Honeypot is easily configured with a single `yml` file. This file is self-validating and, when using [the vcode yaml plugin](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml), practically writes itself.

## **Built to scale**

Honeypot is written entirely in Go which results in a very compact, very useful binary that naturally lends itself to scale. It is written to h-scale quickly and easily without unintended side effects.
