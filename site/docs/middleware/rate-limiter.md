---
tags:
  - middleware
  - rate limiter
---

# 🟢 Rate Limiter

This middleware allows the operator to throttle incoming events from specific sources. It is not intended to be the sole mechanism by which Honeypot is protected from malicious activity, but it does help.

The ratelimiter middleware returns a [429](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429) if the configured threshold is exceeded by a single IP address.


## Configuration

### enabled

  `true` or `false`

### period

One of the following:

`MS` (millisecond)

`S` (second)

`M` (minute)

`H` (hour)

`D` (day)

### limit

**Example:** `10`.

## Sample Configuration

```
rateLimiter:
  enabled: false
  period: S
  limit: 10
```
