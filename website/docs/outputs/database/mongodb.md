---
sidebar_position: 4
---

# Mongodb

**🟢 Supported**

The Mongodb sink writes `valid` and `invalid` events to the configured collections.

Collections are ensured via the nature of Mongodb, so manual creation is not required.

## Sample Mongodb Sink Configuration

```
sinks:
  - name: docsfordays
    type: mongodb
    deliveryRequired: true
    mongoHosts:
      - mongodb1
      - mongodb2
    mongoPort: 27017
    mongoDbName: honeypot
    mongoUser: hpt
    mongoPass: hpt
    validCollection: honeypotValid
    invalidCollection: honeypotInvalid
```
