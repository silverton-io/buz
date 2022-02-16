
# Valid generic event
curl -X POST localhost:8081/gen/p -H 'Content-Type: application/json' -d  '{"contexts": [{"schema": "noschema", "payload": {"something": "else"}}], "payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid"}}}'
# Invalid generic event
curl -X POST localhost:8081/gen/p -H 'Content-Type: application/json' -d  '{"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid", "somethingElse": 10}}}' # Because additional properties not allowed
# Generic event with improperly-specified payload


# Mixture of valid and invalid events, sent to batched event endpoint
curl -X POST localhost:8081/gen/bp -H 'Content-Type: application/json' -d  '[{"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid", "somethingElse": 10}}}, {"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid"}}}, {"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid"}}}, {"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid"}}}, {"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid"}}}]'
# Mixture of valid and invalid due to key mismatch
curl -X POST localhost:8081/gen/bp -H 'Content-Type: application/json' -d  '[{"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid", "somethingElse": 10}}}, {"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid"}}}, {"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid"}}}, {"payload": {"schema": "com.jake/generic/sample/1-0-0.json", "payload": {"id": "myid"}}}, {"payload": {"scheme": "com.jake/generic/sample/1-0-0.json", "data": {"id": "myid"}}}]'
