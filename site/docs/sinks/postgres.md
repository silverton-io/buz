---
tags:
  - sink
  - db
  - postgres
---

# 🟢 Postgres

The Postgres sink writes `valid` and `invalid` events to the configured Postgres tables.

It is especially useful if you already have Postgres running and want to quickly get started with Honeypot-based event tracking.

Tables are ensured upon Honeypot startup, so manual creation is not required.

## Sample Postgres Sink Configuration

```
sinks:
  - name: ol-trusty
    type: postgres
    deliveryRequired: true
    pgHost: localhost
    pgPort: 5432
    pgDbName: honeypot
    pgUser: honeypot
    pgPass: honeypot
    validTable: honeypot_valid
    invalidTable: honeypot_invalid
```
