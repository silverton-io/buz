---
tags:
  - collector
  - input protocol
  - cloud events
---

# 🟢 CloudEvents

## Event Collection Methods

Honeypot listens on a configurable endpoint for incoming `POST` requests of [Cloudevents payloads](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/formats/json-format.md).


This endpoint requires one of the following content types to be designated:

  - `application/cloudevents+json` (for single events)
  - `application/cloudevents-batch+json` (for a batch of events)

**Note!** If a `Content-Type` header is not specified, the event will not be accepted.


## Event Validation Methods

Honeypot validates, annotates, and redirects incoming cloudevents using the [dataschema](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#dataschema) property of incoming events.

**Note!** If the `dataschema` property is not present in incoming events, these events will be redirected to the `invalid` destination(s). This might become configurable in the future, but since Honeypot aims to maintain a high degree of data quality from the point of collection it might not change.


## Sample Cloudevents Configuration

```
cloudevents:
  enabled: true   # Whether or not to enable the Cloudevents collection endpoint
  path: /ce/p     # Path for incoming (single or batch) cloudevents
```