---
tags:
  - sink
  - pubnub
---

# 🟢 PubNub


The PubNub sink writes incoming events to... [PubNub](https://www.pubnub.com/)!


## Sample PubNub Sink Configuration
```
sinks:
  - name: someapp
    type: pubnub
    deliveryRequired: true
    pubnubPubKey: YOUR-PUB-KEY
    pubnubSubKey: YOUR-SUB-KEY
    validChannel: honeypot-valid
    invalidChannel: honeypot-invalid
```
