package main

import (
	"github.com/kishieel/couchdb-query-server-go/pkg/couchgo"
)

func main() {
	server := couchgo.NewQueryServer()
	server.Start()
}

func Update(args couchgo.UpdateInput) couchgo.UpdateOutput {
	var doc couchgo.Document

	if args.Doc == nil {
		doc = couchgo.Document{}
		doc["_id"] = args.Req.Uuid
	} else {
		doc = args.Doc
	}

	doc["data"] = args.Req.Body
	doc["updated"] = true

	return couchgo.UpdateOutput{doc, map[string]any{"body": "Updated"}}
}

func Validate(args couchgo.ValidateInput) couchgo.ValidateOutput {
	return nil
}

// ["reset", {"reduce_limit": true, "timeout": 5000}]
// ["add_fun", "func Map(args couchgo.MapInput) couchgo.MapOutput {\n\tout := make([][2]any, 0)\n\tout = append(out, [2]any{args.Doc[\"_id\"], args.Doc[\"type\"]})\n\n\treturn out\n}"]
// ["add_lib", {"utils": "exports.MAGIC = 42;"}]
// ["map_doc", {"_id": "doc_id", "_rev": "doc_rev", "type": "post", "content": "hello world"}]
// ["map_doc", {"_id": "doc_id", "_rev": "doc_rev", "type": "user", "name": "John Doe"}]
// ["reduce", ["func Reduce(args couchgo.ReduceInput) couchgo.ReduceOutput {\n\tout := 0.0\n\n\tfor _, value := range args.Values {\n\t\tout += value.(float64)\n\t}\n\n\treturn out\n}"], [[[1, "699b"], 10], [[2, "c081"], 20], [[null, "foobar"], 3]]]
// ["rereduce", ["func Reduce(args couchgo.ReduceInput) couchgo.ReduceOutput {\n\tout := 0.0\n\n\tfor _, value := range args.Values {\n\t\tout += value.(float64)\n\t}\n\n\treturn out\n}"], [10, 20, 3]]
// ["ddoc", "new", "_design/myddoc", {"views": {"myview": {"map": "func Map(args couchgo.MapInput) couchgo.MapOutput {\n\tout := make([][2]any, 0)\n\tout = append(out, [2]any{args.Doc[\"_id\"], args.Doc[\"type\"]})\n\n\treturn out\n}"}}}]
// ["ddoc", "_design/myddoc", ["views", "myview", "map"], [{"_id": "doc_id", "_rev": "doc_rev", "type": "post", "content": "hello world"}]]
// ["ddoc", "_design/myddoc", ["filters", "myfilter"], [{"_id": "doc_id", "_rev": "doc_rev", "type": "post"}], {"body": "", "cookie": {}, "form": {}, "headers": {}, "id": "doc_id", "info": {"db_name": "mydb"}, "method": "GET", "path": [], "peer": "", "query": {}, "secobj": {}, "userCtx": {"db": "mydb", "name": "myuser", "roles": []}, "uuid": "uuid"}]
// ["ddoc", "_design/myddoc", ["updates", "myupdate"], [{"_id": "doc_id", "_rev": "doc_rev", "type": "post", "content": "hello world"}, {"body": "", "cookie": {}, "form": {}, "headers": {}, "id": "doc_id", "info": {"db_name": "mydb"}, "method": "GET", "path": [], "peer": "", "query": {}, "secobj": {}, "userCtx": {"db": "mydb", "name": "myuser", "roles": []}, "uuid": "uuid"}]]
// ["ddoc", "_design/myddoc", ["validate_doc_update"], [{"_id": "doc_id"}, {"_id": "doc_id"}, {"name": "myuser", "roles": []}, {"db_name": "mydb"}]]
