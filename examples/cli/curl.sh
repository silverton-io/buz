###########################################################################
# NOTE! These should be turned into tests
###########################################################################


###########################################################################
# Cloudevents
###########################################################################
# Valid event -> valid schema, valid payload
curl -X POST localhost:8080/cloudevents -H 'Content-Type:application/cloudevents+json' -d '{"dataschema":"io.silverton/buz/example/gettingStarted/v1.0.json", "data": {"userId": 10, "name": "jakthom", "action": "didSomething"}}'
# Invalid event -> missing schema, invalid payload
curl -X POST localhost:8080/cloudevents -H 'Content-Type:application/cloudevents+json' -d '{"blah": "blee"}'
# Invalid event -> missing schema, valid payload
curl -X POST localhost:8080/cloudevents -H 'Content-Type:application/cloudevents+json' -d '{"data": {"userId": 10, "name": "jakthom", "action": "didSomething"}}'
# Invalid event -> valid schema, invalid payload (wrong props)
curl -X POST localhost:8080/cloudevents -H 'Content-Type:application/cloudevents+json' -d '{"dataschema":"io.silverton/buz/example/gettingStarted/v1.0.json", "data": {"userId": 10, "name": "jakthom", "activity": "didSomething"}}'
# Invalid event -> valid schema, invalid payload (extra props)
curl -X POST localhost:8080/cloudevents -H 'Content-Type:application/cloudevents+json' -d '{"dataschema":"io.silverton/buz/example/gettingStarted/v1.0.json", "data": {"userId": 10, "name": "jakthom", "action": "didSomething", "somethingElse": "bad"}}'

# Valid event batch -> valid schemas, valid payloads

# Mixed event batch

###########################################################################
# 
###########################################################################