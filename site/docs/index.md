# What is Honeypot?

Honeypot is a multi-protocol event collection, validation, pipelining, and observability system.

It is built to be ***easily-configured***, ***easily-deployed***, and ***easily-maintained***. Yet uncompromising with its speed, guarantees, and operational flexibility. 


## Quickstart

To dive head-first into an example of running Honeypot locally with a three-node [Redpanda](https://github.com/redpanda-data/) cluster, [Kowl](https://github.com/cloudhut/kowl/) for observability, and [Materialize](https://github.com/MaterializeInc/materialize) as a streaming DB please see the [Quickstart](quickstart/getting-started/).

## Philosophy


### Lean on the shoulders of giants
**Data systems and infrastructure are getting very very cool**.

When popular event tracking systems like [Snowplow Analytics](https://github.com/snowplow/snowplow) were first created, deployment systems like K8S and Knative didn't exist. Nor did the "serverless" mindset.


Kafka was an infant. And had not yet created the things that Redpanda solves.


Streaming databases build on the Postgres API were well into the future.


And a data warehouse that would eat the world with its developer-focus yet massively-scalable architecture was yet to be named.


***These are our current-day realities and we want to take advantage of them. All of them.***

### Scale to zero, but also "infinity"


### Validate and redirect on the edge

As event systems 


### Keep operational complexity low


![honeypot](img/honeypot.png)