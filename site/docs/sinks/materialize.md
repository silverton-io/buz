---
tags:
  - sink
  - materialize
---

# 🟢 Materialize

The Materialize sink writes `valid` and `invalid` events to the configured tables.

This sink is especially useful when wanting to try out a streaming database without the overhead of another set of infrastructure.


## Sample Materialize Sink Configuration

```
sinks:
  - name: 🚀🚀🚀
    type: materialize
    deliveryRequired: true
    mzHost: 127.0.0.1
    mzPort: 6875
    mzDbName: materialize
    mzUser: materialize
    mzPass: ""
    validTable: honeypot_valid
    invalidTable: honeypot_invalid
```
