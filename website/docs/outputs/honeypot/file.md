---
sidebar_position: 1
---

# File

**🟢 Supported**

The file sink writes events to respective local files. It is useful, sometimes.

Destination files are ensured on startup, so manual creation is not required.

## Sample File Sink Configuration

```
sinks:
  - name: notgoingfar
    type: file
    deliveryRequired: true
    validFile: honeypot-valid.json
    invalidFile: honeypot-invalid.json
```
