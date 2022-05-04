---
tags:
  - schema cache backend
  - db
  - mongodb
---

# 🟢🎉 MongoDb

The `mongodb` schema cache backend leverages schemas stored in a configurable registry collection.

It is most useful when you want to store `schemas`, `valid events`, and `invalid events` within the same system to reduce infrastructure overhead.

It can be used with any combination of sink(s).


## Sample Mongodb Schema Cache Backend Configuration

```
schemaCache:
  backend:
    type: mongodb
    registryCollection: registry
    mongoHosts:
      - 127.0.0.1
    mongoPort: 27017
    mongoDbName: honeypot
    mongoUser: honeypot
    mongoPass: honeypot
```
