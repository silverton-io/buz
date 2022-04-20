---
tags:
  - schema backend
  - gcs
  - google cloud storage
---

# 🟢 GCS

The `gcs` cache backend uses jsonschemas stored in gcs to back the in-memory schema cache.

## Sample GCS Cache Backend Configuration

```
schemaCache:
  backend:
    type: gcs                 # The backend type
    bucket: honeypot-schemas  # The gcs bucket containing schemas
    path: /                   # The path to consider as root
```
