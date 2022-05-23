---
sidebar_position: 2
---

# Schema Directory

If the `schemaDirectory` is `enabled`, all "currently-cached" schemas are available at the `/schemas` path of the collector.

Individual schemas are also available at `/schemas/$NAME_OF_SCHEMA`.

:::tip For Example
If a schema was named `io.silverton/honeypot/tele/beat/v1.0.json` and the `schemaDirectory` was `enabled` the schema would be available at `$COLLECTOR_URL/schemas/io.silverton/honeypot/tele/beat/v1.0.json`.
:::
