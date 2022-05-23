---
sidebar_position: 2
---

# 🟢 Redpanda/Kafka

The Redpanda/Kafka sink writes `valid` and `invalid` events to the respective topics.


## Sample Redpanda/Kafka Sink Configuration

```
sinks:
  - name: 大熊猫
    type: redpanda
    deliveryRequired: true
    kafkaBrokers:
      - 127.0.0.1:9092
    validTopic: honeypot-valid
    invalidTopic: honeypot-invalid
```
