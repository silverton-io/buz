---
sidebar_position: 6
---

# 🟢 HTTP/S

The http/s sink writes events via batched `POST` requests to the configured `valid` and `invalid` urls. Without discretion.


## Sample HTTP/S Sink Configuration

```
sinks:
  - name: somewheres
    type: https
    deliveryRequired: true
    validUrl: https://your-endpoint.net/valid-events-here
    invalidUrl: https://your-endpoint.net/invalid-events-here
```
