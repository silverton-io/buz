---
sidebar_position: 5
---

# 🟢 AWS Kinesis Firehose

The Kinesis Firehose sink writes `valid` and `invalid` events to the configured streams. It is especially useful when wanting to write incoming events directly to S3.

## Sample AWS Kinesis Firehose Sink Configuration

```
sinks:
  - name: straightshots3
    type: kinesis-firehose
    deliveryRequired: true
    validStream: honeypot-valid
    invalidStream: honeypot-invalid
```
