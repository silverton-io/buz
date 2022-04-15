---
tags:
  - collector
  - input protocol
  - self describing
---

# 🟢 Self-Describing Events


## Sample Self-Describing Event Configuration
```
  generic:
    enabled: true         # Whether or not to enable generic self-describing events
    path: /gen/p          # The path to use for POST-based SD event tracking
    contexts: 
      rootKey: contexts   # The contexts root key (contexts)
      schemaKey: schema   # The contexts schema key (contexts.schema)
      dataKey: data       # The contexts data key (contexts.data)
    payload:
      rootKey: payload    # The payload root key (payload)
      schemaKey: schema   # The payload schema key (payload.schema)
      dataKey: data       # The payload data key (payload.data)
```
