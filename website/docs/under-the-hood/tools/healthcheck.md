---
sidebar_position: 4
---

# Healthcheck

Honeypot ships with a `/healthcheck` endpoint by default.

This endpoint can be both`disabled` altogether or pointed at a different path.

Configuration is as follows:

```
app:
  health:
    enabled: true
    path: /health
```
