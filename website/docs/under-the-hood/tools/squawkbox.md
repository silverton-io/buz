---
sidebar_position: 3
---

# Squawkbox


The `squawbox` is a simple mechanism for quick feedback - incoming events are validated and enveloped, and the envelope is returned as response body.


It is most helpful when building out a tracking implementation, doing local development, etc.


`Squawkbox` is configurable - it can be disabled altogether or shifted to a set of alternative paths:

```
squawkBox:
  enabled: true
  cloudeventsPath: /sqwk/ce
  snowplowPath: /sqwk/sp
  genericPath: /sqwk/gen
```
