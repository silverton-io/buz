---
tags:
  - schema cache backend
  - db
  - materialize
---

# 🟢🎉 Materialize

The `materialize` schema cache backend leverages schemas stored in a configurable registry table.

It is most useful when you want to store `schemas`, `valid events`, and `invalid events` within the same system to reduce infrastructure overhead.

It can be used with any combination of sink(s).


## Sample Materialize Schema Cache Backend Configuration

```
schemaCache:
  backend:
    type: materialize
    registryTable: registry
    mzHost: localhost
    mzPort: 6875
    mzDbName: materialize
    mzUser: materialize
    mzPass: ""
```
