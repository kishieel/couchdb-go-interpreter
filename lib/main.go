package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	server := NewQueryServer()

	for scanner.Scan() {
		var tmp []interface{}

		if err := json.Unmarshal(scanner.Bytes(), &tmp); err != nil {
			Respond([]string{"error", "unnamed_error", err.Error()})
		}

		kind := GetCommandKind(tmp...)

		factory, found := CommandRegistry[kind]
		if !found {
			Respond([]string{"error", "unknown_command", fmt.Sprintf("Unknown command type: %s", kind)})
			continue
		}

		command := factory()
		command.Parse(tmp...)

		dispatcher, found := CommandDispatcher[kind]
		if !found {
			Respond([]string{"error", "unknown_command", fmt.Sprintf("Unknown command type: %s", kind)})
			continue
		}

		dispatcher(server, command)
	}
}

// ["reset", {"reduce_limit": true, "timeout": 5000}]
// ["add_fun", "func Map(args couchgo.MapInput) couchgo.MapOutput {\n\tout := make([][2]any, 0)\n\tout = append(out, [2]any{args.Doc[\"_id\"], args.Doc[\"type\"]})\n\n\treturn out\n}"]
// ["add_lib", {"utils": "exports.MAGIC = 42;"}]
// ["map_doc", {"_id": "doc_id", "_rev": "doc_rev", "type": "post", "content": "hello world"}]
// ["map_doc", {"_id": "doc_id", "_rev": "doc_rev", "type": "user", "name": "John Doe"}]
// ["reduce", ["func Reduce(args couchgo.ReduceInput) couchgo.ReduceOutput {\n\tout := 0.0\n\n\tfor _, value := range args.Values {\n\t\tout += value.(float64)\n\t}\n\n\treturn out\n}"], [[[1, "699b"], 10], [[2, "c081"], 20], [[null, "foobar"], 3]]]
// ["rereduce", ["func Reduce(args couchgo.ReduceInput) couchgo.ReduceOutput {\n\tout := 0.0\n\n\tfor _, value := range args.Values {\n\t\tout += value.(float64)\n\t}\n\n\treturn out\n}"], [10, 20, 3]]
// ["ddoc", "new", "_design/myddoc", {"views": {"myview": {"map": "function(doc) { emit(doc._id, doc); }"}}}]
// ["ddoc", "_design/myddoc", ["views", "myview", "map"], [{"_id": "doc_id", "_rev": "doc_rev", "type": "post", "content": "hello world"}]]
// ["ddoc", "_design/myddoc", ["filters", "myfilter"], [{"_id": "doc_id", "_rev": "doc_rev", "type": "post"}], {"body": "", "cookie": {}, "form": {}, "headers": {}, "id": "doc_id", "info": {"db_name": "mydb"}, "method": "GET", "path": [], "peer": "", "query": {}, "secobj": {}, "userCtx": {"db": "mydb", "name": "myuser", "roles": []}, "uuid": "uuid"}]
// ["ddoc", "_design/myddoc", ["updates", "myupdate"], [{"_id": "doc_id", "_rev": "doc_rev", "type": "post", "content": "hello world"}, {"body": "", "cookie": {}, "form": {}, "headers": {}, "id": "doc_id", "info": {"db_name": "mydb"}, "method": "GET", "path": [], "peer": "", "query": {}, "secobj": {}, "userCtx": {"db": "mydb", "name": "myuser", "roles": []}, "uuid": "uuid"}]]
// ["ddoc", "_design/myddoc", ["validate_doc_update"], [{"_id": "doc_id"}, {"_id": "doc_id"}, {"name": "myuser", "roles": []}, {"db_name": "mydb"}]]
