# What is Honeypot?

## Event Streaming for The Rest of Us

Honeypot is a multi-protocol event collection, validation, pipelining, and observability system.

It is designed to be ***easily-configured***, ***easily-deployed***, and ***easily-maintained***.

Yet uncompromising with its speed, guarantees, and operational flexibility. 


# Quickstart

To dive head-first into an example of running Honeypot locally with a three-node [Redpanda](https://github.com/redpanda-data/) cluster, [Kowl](https://github.com/cloudhut/kowl/), and [Materialize](https://github.com/MaterializeInc/materialize) please see the [Quickstart](quickstart/getting-started/).

# Philosophy


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


![honeypot](img/honeypot.png)