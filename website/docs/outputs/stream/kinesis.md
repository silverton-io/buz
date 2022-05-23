---
sidebar_position: 3
---


# AWS Kinesis

**🟢 Supported**

The Kinesis sink writes `valid` and `invalid` events to the configured streams.

## Sample AWS Kinesis Sink Configuration

```
sinks:
  - name: zoom
    type: kinesis
    deliveryRequired: true
    validStream: honeypot-valid
    invalidStream: honeypot-invalid
```
