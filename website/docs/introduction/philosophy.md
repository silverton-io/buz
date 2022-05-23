---
sidebar_position: 3
---

# Philosophy


## Build new systems on proven API's and mental models

There are some very good API's out there. Oftentimes these API's were originally built using best-in-class tech, but said tech has since been supplanted. Why recreate the world when advancing an existing API or mental model will do the trick?

Examples of building new technology on top of pre-existing API's include:

- [Redpanda](https://redpanda.com/),  which is built on [Kafka](https://kafka.apache.org/documentation/)'s API.

- [Materialize](https://materialize.com/), [Timescale](https://www.timescale.com/), [CockroachDB](https://www.cockroachlabs.com/product/), and [PipelineDB](https://www.confluent.io/blog/pipelinedb-team-joins-confluent/) which are all built on the Postgres wire protocol.

- [Airbyte](https://airbyte.com/), which builds upon the conceptual model of Fivetran or Stitch (and the technical model of [Singer](https://www.singer.io/)).


## Stand on the shoulders of giants
**Data systems and infrastructure are getting very very cool**.

When popular event tracking systems like [Snowplow Analytics](https://github.com/snowplow/snowplow) were first created, deployment systems like K8S and Knative didn't exist. Nor did the "serverless" mindset.


Kafka was an infant. And had not yet created the things that Redpanda solves.


Streaming databases built on the Postgres API were years into the future.


And a data warehouse that would eat the world with its developer-focus yet massively-scalable architecture was yet to be named.


***These are all present-day realities and we want to build upon them with an eye to the future.***

## Scale to zero, but also "infinity"

Why pay for what you don't use? Or completely rearchitect systems as volume grows or demands change?

Serverless scales to zero, and then back up again...

Snowflake scales to zero, and then back up again...

***Event collection systems should too.***


## Validate and redirect on the edge

Data should be validated and redirected as soon as it enters collection infrastructure, not near the end of the process (or in the data warehouse). The faster data is declared "valid" or "invalid", the faster it can be used and acted upon.


## Keep operational complexity low

And last but certainly not least, engineers should be able to maintain and advance event collection efforts without complexity or cost explosion.


# Would you like to know more?

If you would like to know more or follow the project, **[check out the roadmap](/introduction/roadmap)** or sign up for **[Insiders-Only Updates](/insiders-only)**.
