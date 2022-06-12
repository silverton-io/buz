---
sidebar_position: 3
---

# Indicative

**🟢 Supported**

This sink writes formatted envelopes to Indicative for easy visualization and analysis.

:::tip Note!
All nested payload properties are flatted and dot-separated, so this:
```
{
    "users": {
        "jane": {
            "age": 40
        }
    }
}
```

Will show up in Indicative as a property of `users.jane.age` having the value of `40`.
:::


All associated envelope properties are passed along to Indicative.


## Sample Indicative Sink Configuration

```
sinks:
  - name: sweetCharts
    type: indicative
    deliveryRequired: true
    indicativeApiKey: YOUR-API-KEY-HERE
```
