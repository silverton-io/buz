---
sidebar_position: 10
---


# 🟢 MySQL

The MySQL sink writes `valid` and `invalid` events to the configured tables.

It is especially useful if you already have MySQL running and want to quickly get started with Honeypot-based event tracking.

Tables are ensured upon Honeypot startup, so manual creation is not required.

## Sample MySQL Sink Configuration

```
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

