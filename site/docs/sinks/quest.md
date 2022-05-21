---
tags:
  - sink
  - db
  - quest
---

# 🟢🎉 Quest

The Quest sink writes `valid` and `invalid` events to the configured tables.

It is especially useful if you already have Quest running and want to quickly get started with Honeypot-based event tracking.

Tables are ensured upon Honeypot startup, so manual creation is not required.

## Sample Quest Sink Configuration

```
sinks:
  - name: ts-rocket
    type: quest
    deliveryRequired: true
    questHost: localhost
    questPort: 5432
    questDbName: honeypot
    questUser: admin
    questPass: quest
    validTable: honeypot_valid
    invalidTable: honeypot_invalid
```
