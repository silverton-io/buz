---
tags:
  - sink
  - db
  - elasticsearch
---

# 🟢 Elasticsearch

The Elasticsearch sink loads `valid` and `invalid` events into the configured indices.

Indices are ensured via the nature of elasticsearch, so manual creation is not required.

## Sample Elasticsearch Sink Configuration

```
sinks:
  - name: loggin
    type: elasticsearch
    deliveryRequired: true
    elasticsearchHosts: 
      - http://es1:9200
    elasticsearchUsername: elastic
    elasticsearchPassword: elastic
    validIndex: honeypot-valid
    invalidIndex: honeypot-invalid
```
