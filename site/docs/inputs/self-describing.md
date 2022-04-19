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
  $CONTEXTS_KEY: [
    {$CONTEXT_SCHEMA_KEY: "some-context-key", $CONTEXT_DATA_KEY: {"context-data": "here"}},
    {$CONTEXT_SCHEMA_KEY: "another-context-key", $CONTEXT_DATA_KEY: {"more-context-data": "here"}},
  ],
  $PAYLOAD_KEY: {
    $SCHEMA_KEY: "some-key",
    $DATA_KEY: {"data": "here"}
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
      schemaKey: schema   # The contexts schema key (contexts.schema)
      dataKey: data       # The contexts data key (contexts.data)
    payload:
      rootKey: payload    # The payload root key (payload)
      schemaKey: schema   # The payload schema key (payload.schema)
      dataKey: data       # The payload data key (payload.data)
```
