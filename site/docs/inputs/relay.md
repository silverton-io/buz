---
tags:
  - collector
  - input protocol
  - relay
---

# 🟢 Honeypot Relay


## Collection Methods

Honeypot is capable of collecting and distributing events relayed from other honeypot instances, which allows for operational flexibility in a set of use cases.


## Validation Method

Relayed messages **are not re-validated** since messages are validated at point of initial collection.


## Sample Relay Configuration

```
relay:
  enabled: true     # Whether or not to accept relayed events
  path: /relay      # Path for incoming relayed events
```
