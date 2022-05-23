---
sidebar_position: 2
---

# 🟢 Snowplow


## Protocols

At the end of the day Snowplow leverages *two* event protocols - the original [Tracker Protocol](https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/snowplow-tracker-protocol/) and [Custom Self-Describing Events](https://docs.snowplowanalytics.com/docs/collecting-data/collecting-from-own-applications/javascript-trackers/javascript-tracker/javascript-tracker-v2/tracking-specific-events/#tracking-custom-self-describing-events).

Honeypot supports both, but does so in a way that seamlessly blends the data model of Self-Describing Events and the traditional Tracker Protocol. It also *validates and redirects tracker-protocol events in the same manner as self-describing events.*

## Collection Methods

Snowplow uses two HTTP verbs for event collection:

  - `GET` (request consisting of a single event, as defined via query params)
  - `POST` (request consisting of event batches, as defined via json payloads)

It leverages three primary event collection endpoints:

  - `/i` (pixel endpoint used for `GET`-based tracking)
  - `r/tp2` (redirect endpoint used for `GET`-based tracking)
  - `com.snowplow.analytics.snowplow/tp2` (endpoint used for `POST`-based tracking)

Honeypot supports all of the above, but allows various functionality to be enabled/disabled as needed.

## Validation Method

Honeypot uses a set of static schemas to validate incoming "tracker protocol" events, and the associated `schema` property to validate incoming "self describing" events.

## Navigating Ad Blockers

Honeypot supports configurable collection endpoints so tracking does not get blocked by ever-expanding ad blocker lists.


## Sample Snowplow Configuration

```
inputs:
  snowplow:
    enabled: true               # Whether or not to enable Snowplow event collection
    standardRoutesEnabled: true # Whether or not to enable Snowplow's standard routes
    openRedirectsEnabled: true  # Whether or not to enable open redirects
    getPath: /plw/g             # The custom path for get-based tracking
    postPath: /plw/p            # The custom path for post-based tracking
    redirectPath: /plw/r        # The custom path for open redirect tracking
    anonymize:
      ip: false                 # Whether or not to anonymize ip's
      userId: false             # Whether or not to anonymize userid's
