---
sidebar_position: 2
---

# Rate Limiter

**🟢 Supported**

The rate limiter middleware allows the operator to throttle incoming events from specific sources. It is not intended to be the sole mechanism by which Honeypot is protected from malicious activity, but it does help.


:::danger Status Code
A [429](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429) is returned if configured threshold is exceeded by a single IP address.
:::


## Sample Rate Limiter Middleware Configuration

```
rateLimiter:
  enabled: false
  period: S
  limit: 10
```
