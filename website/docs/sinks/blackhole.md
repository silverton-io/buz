---
sidebar_position: 14
---

# 🟢 Blackhole


The blackhole sink is the equivalent of sinking events to `/dev/null`.

It is primarily useful as a development tool or when collecting events in non-production environments if you don't want to sink them anywhere.


## Sample Blackhole Sink Configuration

```
sinks:
  - name: supermassive
    type: blackhole
    deliveryRequired: true
```
