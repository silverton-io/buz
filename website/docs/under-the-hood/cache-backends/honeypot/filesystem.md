# Filesystem

**🟢 Supported**

The `fs` cache backend uses jsonschemas stored on the local filesystem to back the in-memory schema cache.

## Sample Filesystem Schema Cache Backend Configuration

```
schemaCache:
  backend:
    type: fs                        # The backend type
    path: /some/path/somewhere      # The path to consider as root
```
