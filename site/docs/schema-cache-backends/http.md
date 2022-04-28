---
tags:
  - schema backend
  - http
  - https
---

# 🟢 HTTP/S

The `http` and `https` cache backends use jsonschemas stored at remote HTTP paths to back the in-memory schema cache.

## Sample Filesystem Cache Backend Configuration

```
schemaCache:
  backend:
    type: https                     # The backend type
    host: registry.silverton.io     # The schema host
    path: /some/path/somewhere      # The path to consider as root
```
