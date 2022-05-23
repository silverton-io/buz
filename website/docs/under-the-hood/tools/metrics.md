---
sidebar_position: 2
---

# Event-Level Metrics

Each Honeypot instance has an onboard metrics endpoint located (by default) at `/stats`.

Stats are representative of aggregate `invalid` and `valid` event volume, at an event level, since startup.


The stats endpoint is configurable - it can be disabled altogether or located at a different path:
```
app:
  stats:
    enabled: true
    path: /stats
```