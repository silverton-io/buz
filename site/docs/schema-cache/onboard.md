---
tags:
  - schema cache
  - onboard
---

# 🟢 Onboard


## Sample Schema Cache Configuration
```
schemaCache:
  schemaCacheBackend:         # Cache backend configuration
    type: fs                  # The type of remote schema cache backend
    path: ./schemas/          # The root path of a schema cache backend of type "fs"
  ttlSeconds: 300             # The number of seconds to keep a schema in the onboard cache
  maxSizeBytes: 104857600     # The max size, in bytes, of the onboard cache
  purge:
    enabled: true             # Whether or not to enable a cache purge route
    path: /c/purge            # The path of the cache purge route
  schemaDirectory:
    enabled: true             # Whether or not to enable schema directory routes
```
