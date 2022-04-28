---
tags:
  - sink
  - file
---

# 🟢 File

The file sink writes events to respective local files. It is useful, sometimes.


## Sample File Sink Configuration

```
sinks:
  - name: notgoingfar
    type: file
    deliveryRequired: true
    validFile: honeypot-valid.json
    invalidFile: honeypot-invalid.json
```
