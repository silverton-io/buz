---
tags:
  - sink
  - db
  - clickhouse
---

# 🟢 Clickhouse

The Clickhouse sink loads `valid` and `invalid` events into the configured tables.

Table existence is ensured each time Honeypot starts up so manual creation is not required.

## Sample Clickhouse Sink Configuration

```
sinks:
  - name: houseofclicks
    type: clickhouse
    deliveryRequired: true
    clickhouseHost: 127.0.0.1
    clickhousePort: 9000
    clickhouseDbName: honeypot
    clickhouseUser: honeypot
    clickhousePass: honeypot
    validTable: honeypot_valid
    invalidTable: honeypot_invalid
```
