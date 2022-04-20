---
tags:
  - schema cache
  - onboard
---

# 🟢 Onboard

In order to drastically improve event validation speed, schemas are cached onboard each running honeypot for the configure period of time.


## Schema Directory

If the `schemaDirectory` is `enabled`, all "currently-cached" schemas are available at the `/schemas` path of the collector.

Individual schemas are also available at `/schemas/$NAME/$OF/$SCHEMA`.

For example, if a particular schema had a name of `com.silverton.io/honeypot/tele/beat/v1.0.json` it would be available at `/schemas/com.silverton.io/honeypot/tele/beat/v1.0.json`.


## Sample Schema Cache Configuration
```
schemaCache:
  ttlSeconds: 300             # The number of seconds to keep a schema in the onboard cache
  maxSizeBytes: 104857600     # The max size, in bytes, of the onboard cache
  purge:
    enabled: true             # Whether or not to enable a cache purge route
    path: /c/purge            # The path of the cache purge route
  schemaDirectory:
    enabled: true             # Whether or not to enable schema directory routes
```
