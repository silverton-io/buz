---
tags:
  - sink
  - mysql
---

# 🟢 MySQL

The MySQL sink writes `valid` and `invalid` events to the configured tables.

It is especially useful if you already have MySQL running and want to quickly get started with Honeypot-based event tracking.

## Sample MySQL Sink Configuration

```
sinks:
sinks:
  - name: whoa-nelly
    type: mysql
    deliveryRequired: true
    mysqlHost: localhost
    mysqlDbName: honeypot
    mysqlPort: 3306
    mysqlUser: honeypot
    mysqlPass: honeypot
    validTable: honeypot_valid
    invalidTable: honeypot_invalid
```
