---
tags:
  - middleware
  - request logger
---

# 🟢 Request Logger

The request logger middleware does exactly what it sounds like!

It logs requests in the following form:

```
{"level":"info","request":{"responseCode":200,"requestDuration":7412000,"requestDurationForHumans":"7.412ms","clientIp":"127.0.0.1:59221","requestMethod":"POST","requestUri":"/com.snowplowanalytics.snowplow/tp2"},"time":"2022-04-28T00:58:15-04:00"}
```


## Sample Request Logger Configuration

```
requestLogger:
  enabled: true
```
