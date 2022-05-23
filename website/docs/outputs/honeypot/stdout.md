---
sidebar_position: 4
---


# Stdout

**🟢 Supported**


The stdout sink writes colorized events to.... stdout! It is especially useful when wanting feedback during development or when taking Honeypot for a test drive.


## Sample Stdout Sink Configuration

```
sinks:
  - name: console
    type: stdout
    deliveryRequired: true
```
