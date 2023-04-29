package main

import (
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"log"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func main() {
	c := jsonschema.NewCompiler()
	c.AssertContent = true
	c.Decoders["hex"] = hex.DecodeString
	c.MediaTypes["application/xml"] = func(b []byte) error {
		return xml.Unmarshal(b, new(interface{}))
	}

	schema := `{ "$id": "io.silverton/buz/example/gettingStarted/v1.0.json", "title": "io.silverton/buz/example/gettingStarted/v1.0.json", "description": "A getting started event used to bootstrap and demonstrate validation", "type": "object", "properties": { "userId": { "type": "integer", "description": "The id of the user" }, "name": { "type": "string", "description": "The name of the user" }, "action": { "type": "string", "description": "The associated user action" } }, "additionalProperties": false, "required": [ "userId", "name", "action" ] }`
	instance := `{"action":"didSomething","name":"jakthom","userId":"10"}`

	if err := c.AddResource("schema.json", strings.NewReader(schema)); err != nil {
		log.Fatalf("%v", err)
	}

	sch, err := c.Compile("schema.json")
	if err != nil {
		log.Fatalf("%#v", err)
	}

	var v interface{}
	if err := json.Unmarshal([]byte(instance), &v); err != nil {
		log.Fatal(err)
	}

	if err = sch.Validate(v); err != nil {
		log.Fatalf("%#v", err)
	}
	// Output:
}
