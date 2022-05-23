# Clickhouse

**🟢 Supported**

The `clickhouse` schema cache backend leverages schemas stored in a configurable registry table.

It is most useful when you want to store `schemas`, `valid events`, and `invalid events` within the same system to reduce infrastructure overhead.

It can be used with any combination of sink(s).

## Sample Clickhouse Schema Cache Backend Configuration

```
schemaCache:
  backend:
    type: clickhouse
    registryTable: registry
    clickhouseHost: 127.0.0.1
    clickhousePort: 9000
    clickhouseDbName: honeypot
    clickhouseUser: honeypot
    clickhousePass: honeypot
```
