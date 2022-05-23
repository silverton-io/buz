---
tags:
  - sink
  - message broker
  - nats
---

# 🟢🎉 NATS Jetstream

The nats sink writes `valid` and `invalid` events to the configured streams.

It is especially useful if you already have Jetstream running and want to quickly get started with Honeypot-based event tracking.


## Sample NATS Sink Configuration

```
sinks:
  - name: streamjet
    type: nats
    deliveryRequired: true
    natsHost: nats
    natsUser: someuser
    natsPass: somepass
    validStream: honeypot.valid
    invalidStream: honeypot.invalid
```
