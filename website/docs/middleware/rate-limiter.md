---
sidebar_position: 3
---

# 🟢 Rate Limiter

The rate limiter middleware allows the operator to throttle incoming events from specific sources. It is not intended to be the sole mechanism by which Honeypot is protected from malicious activity, but it does help.

The ratelimiter middleware returns a [429](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429) if the configured threshold is exceeded by a single IP address.


## Sample Rate Limiter Middleware Configuration

```
rateLimiter:
  enabled: false
  period: S
  limit: 10
```
