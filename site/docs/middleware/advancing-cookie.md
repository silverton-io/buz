---
tags:
  - middleware
  - advancing cookie
---

# 🟢 Advancing Cookie

The advancing cookie middleware sets a server-side identity cookie when enabled. It is used to track across authentication boundaries, or to roll up events and activity sessions to a single user regardless of auth status.

## Configuration

### enabled

Whether or not to activate the advancing cookie middleware.

**Example:** `true` or `false`

### name

The name of the cookie.

**Example:** `nuid`

### secure

The secure cookie attribute.

**Example:** `true` or `false`

### ttlDays

The number of days to persist the advancing cookie.

**Example:** `365`

### domain

The domain to persist the advancing cookie on.

**Example:** `silverton.io`

### path

The path to persist the advancing cookie on.

**Example:** `/`

### sameSite

The `sameSite` attribute of the cookie

**Example:** `Strict`, `Lax`, `None`

## Sample Configuration


```
cookie:
  enabled: true
  name: nuid
  secure: true
  ttlDays: 365
  domain: localhost
  path: /
  sameSite: Lax
```
