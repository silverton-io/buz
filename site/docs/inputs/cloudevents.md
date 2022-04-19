---
tags:
  - collector
  - input protocol
  - cloud events
---

# 🟢 CloudEvents

## Collection Method

Honeypot listens on a configurable endpoint for incoming `POST` requests of [Cloudevents payloads](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/formats/json-format.md).


This endpoint requires one of the following content types to be designated:

  - `application/cloudevents+json` (for single cloudevents)
  - `application/cloudevents-batch+json` (for a batch of cloudevents)

**Note!** If a `Content-Type` header is not specified, the event will not be accepted.


## Validation Method

Honeypot validates incoming cloudevents using the [dataschema](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#dataschema) property of each event.

**Note!** If the `dataschema` property is not present in incoming events, these events will be redirected to the `invalid` destination(s).


## Sample Cloudevents Configuration

```
cloudevents:
  enabled: true   # Whether or not to enable the Cloudevents collection endpoint
  path: /ce/p     # Path for incoming (single or batch) cloudevents
```
