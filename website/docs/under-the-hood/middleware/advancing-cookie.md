---
sidebar_position: 5
---

# Advancing Cookie

**🟢 Supported**

The advancing cookie middleware sets a server-side identity cookie when enabled. It is used to track across authentication boundaries, or to roll up events and activity sessions to a single user regardless of auth status.


## Sample Advancing Cookie Middleware Configuration


```
cookie:
  enabled: true
  name: nuid
  secure: true
  ttlDays: 365
  domain: localhost
  path: /
  sameSite: Lax
```
