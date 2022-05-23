---
sidebar_position: 3
---


# Pixel

**🟢 Supported**

## Collection Method

### Parameter Payloads

Honeypot supports collecting payloads via url query params.

While this method of data collection is useful it does have drawbacks such as [max uri lengths](https://stackoverflow.com/questions/812925/what-is-the-maximum-possible-length-of-a-query-string) and the inability to explicitly declare parameter types.

The *good* thing about the pixel input is it is extremely simple to get started with.

**For example** -> if `/pxl` is configured as a pixel input, submitting a `GET` request to `/pxl/?hello=world&userId=10` will send a payload of `{"hello": "world", "userId": "10"}` to the configured sinks. No sdk's necessary.

### Base64 Encoded Parameter Payloads

Honeypot supports a "special" query param, `hbp`, by which b64 encoded payloads can be collected.

**For example** -> if `/pxl` is configured as a pixel input, submitting a `GET` request to `/pxl?hbp=eyJoZWxsbyI6IndvcmxkIn0` will send a payload of `{"hello":"world"}` to the configured sinks.


## Validation Method

:::danger Please Note
Honeypot **does not yet validate** incoming pixel-based payloads - they are assumed "ok" and are sunk to the "good" sink(s).
:::


## Sample Pixel Configuration

```
pixel:
  enabled: true
  paths:
    - name: default
      path: /pxl/d
    - name: secondary
      path: /pxl/scnd
```
