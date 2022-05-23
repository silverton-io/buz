---
sidebar_position: 4
---


# 🟢 Webhook


## Collection Method

Honeypot is capable of collecting both **named** and **unnamed** webhooks. This is designed so webhooks from multiple sources can be fired into a single set of collection infrastructure, and the sources will retain distinct identifiers.

### Named webhooks

Named webhooks are the recommended approach and are relatively simple - the url path following what is specified in the `webhook` config block becomes the identifier for all webhooks sent there.


**For example:** you have configured the webhook path to be `/pooh-bear` via:

```
webhook:
  path: /pooh-bear
```


Any webhook payloads fired to `/pooh-bear/revenue/stripe/v1` will be identified as `revenue/stripe/v1`, payloads fired to `/pooh-bear/gitlab` will be identified as `gitlab`, payloads fired to `/pooh-bear/d016bb00-02db-4cc6-9852-e29c7cf3aa57` will be identified as `d016bb00-02db-4cc6-9852-e29c7cf3aa57`, etc....



### Unnamed webhooks

Unnamed webhooks are the fallback, but this functionality comes with the inevitability of creating a pile of random json with limited context.

**Which you probably don't want! *So use at your own risk.***

Unnamed webhooks are also relatively straight-forward - the configured url path in the `webhook` config block acts as an unnamed catch-all.

**For example:** you have configured the webhook path to be `/christopher-robbin` via:

```
webhook:
  path: /christopher-robbin
```

Any webhook payloads fired to `/christopher-robbin` will be collected and passed along as "unnamed webhooks".

## Validation Method

Honeypot does **not** validate incoming webhooks - they are assumed "ok".


## Sample Webhook Configuration

```
webhook:
  enabled: true     # Whether or not to enable webhook
  path: /wb/hk      # Path for incoming webhooks
```
