###########################################################################
# These should be turned into tests
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
curl -X POST localhost:8080/cloudevents -H 'Content-Type:application/cloudevents+json' -d '[{"dataschema":"io.silverton/buz/example/gettingStarted/v1.0.json", "data": {"userId": 10, "name": "jakthom", "action": "didSomething"}}, {"dataschema":"io.silverton/buz/example/gettingStarted/v1.0.json", "data": {"userId": 10, "name": "jakthom", "action": "didSomething"}}]'
# Mixed event batch
curl -X POST localhost:8080/cloudevents -H 'Content-Type:application/cloudevents+json' -d '[{"dataschema":"io.silverton/buz/example/gettingStarted/v1.0.json", "data": {"userId": 10, "name": "jakthom", "action": "didSomething", "somethingElse": "bad"}}, {"dataschema":"io.silverton/buz/example/gettingStarted/v1.0.json", "data": {"userId": 10, "name": "jakthom", "action": "didSomething"}}]'

###########################################################################
# Pixel
###########################################################################
# Valid arbitrary event -> arbitrary schema, valid payload
curl -X GET "localhost:8080/pixel?id=10"
# Valid arbitrary event -> explicit schema, valid payload
curl -X GET "localhost:8080/pixel/io.silverton/buz/example/generic/sample/v1.0?id=10"
# Invalid arbitrary event -> explicit schema, invalid payload (want int, got string)
curl -X GET "localhost:8080/pixel/io.silverton/buz/example/gettingStarted/v1.0?userId=10&name=george&action=clicked"

###########################################################################
# Webhook
###########################################################################
# Valid arbitrary event -> arbitrary schema, valid payload

# Valid arbitrary event -> explicit schema, valid payload

# Invalid arbitrary event -> explicit schema, invalid payload
