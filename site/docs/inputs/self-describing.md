---
tags:
  - collector
  - input protocol
  - self describing
---

# 🟢 Self-Describing Events

## Collection Method
Honeypot listens on a configurable endpoint for incoming `POST` requests of `self describing` payloads, structured as:

```
{
  $CONTEXTS_ROOT_KEY: {
    "some-context-schema": {"context-data": "here"},
    "another-context-schema": {"more-context-data": "here"}
  },
  $PAYLOAD_ROOT_KEY: {
    $PAYLOAD_SCHEMA_KEY: "some-key",
    $PAYLOAD_DATA_KEY: {"data": "here"}
  }
}
```

This (configured by you!) endpoint accepts **batches of self-describing events** and **single self-describing events**. It requires a `Content-Type` header of `application/json`.

**Note!** If a `Content-Type` header is not specified, the event will not be accepted.


## Validation Method

Honeypot uses the schema defined at `$PAYLOAD_KEY.$SCHEMA_KEY` to validate each payload.


## Sample Self-Describing Event Configuration
```
  generic:
    enabled: true         # Whether or not to enable generic self-describing events
    path: /gen/p          # Path for incoming self-describing events
    contexts: 
      rootKey: contexts   # The contexts root key (contexts)
    payload:
      rootKey: payload    # The payload root key (payload)
      schemaKey: schema   # The payload schema key (payload.schema)
      dataKey: data       # The payload data key (payload.data)
```
