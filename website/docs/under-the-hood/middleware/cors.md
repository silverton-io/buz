---
sidebar_position: 1
---

# CORS

**🟢 Supported**

The cors middleware ensures Honeypot is able to track events across a set of domains.

It allows the following headers to be entirely configurable:

* `Access-Control-Allow-Origin`
* `Access-Control-Allow-Credentials`
* `Access-Control-Allow-Headers`
* `Access-Control-Allow-Methods`
* `Access-Control-Max-Age`


## Sample CORS Middleware Configuration

```
cors:
  enabled: true
  allowOrigin:
    - "*"
  allowCredentials: true
  allowMethods:
    - POST
    - OPTIONS
    - GET
  maxAge: 86400
```
