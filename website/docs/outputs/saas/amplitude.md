---
sidebar_position: 2
---

# Amplitude

**🟢 Supported**

This sink writes formatted envelopes to Amplitude for easy visualization and analysis.

:::tip Note!
All nested payload properties are flatted and dot-separated, so this:
```
{
    "topLevel": {
        "secondLevel": {
            "thirdLevel": 10
        }
    }
}
```

Will show up in Amplitude as a property of `topLevel.secondLevel.thirdLevel`, having the value of `10`.
:::


All associated envelope properties are passed along to Amplitude.


## Sample Amplitude Sink Configuration

```
sinks:
  - name: sweetCharts
    type: amplitude
    deliveryRequired: true
    amplitudeApiKey: YOUR-API-KEY-HERE
    amplitudeRegion: standard # Either standard or eu
```
