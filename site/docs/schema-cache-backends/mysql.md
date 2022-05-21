---
tags:
  - schema cache backend
  - db
  - mysql
---

# 🟢 MySQL

The `mysql` schema cache backend leverages schemas stored in a configurable registry table.

It is most useful when you want to store `schemas`, `valid events`, and `invalid events` within the same system to reduce infrastructure overhead.

It can be used with any combination of sink(s).


## Sample Mysql Schema Cache Backend Configuration

```
schemaCache:
  backend:
    type: mysql
    registryTable: registry
    mysqlHost: localhost
    mysqlPort: 3306
    mysqlDbName: honeypot
    mysqlUser: honeypot
    mysqlPass: honeypot
```
