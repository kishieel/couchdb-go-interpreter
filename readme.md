# CouchDB Query Server Go

The CouchDB Query Server Go is a Go implementation of the CouchDB Query Server protocol for CouchDB 3.3.0 or later.
It is a drop-in replacement for the JavaScript implementation that ships with CouchDB.

## Usage

To use the CouchDB Query Server Go, you need to add `couchgo` executable to the CouchDB and add the following environment variable:

```shell
export COUCHDB_QUERY_SERVER_GO=/path/to/couchgo
```

Then you can write your design documents in Go and use them in CouchDB.

```json
{
  "_id": "_design/example",
  "views": {
    "view-1": {
      "map": "func Map(args couchgo.MapInput) couchgo.MapOutput {\n\tout := make([][2]any, 0)\n\tout = append(out, [2]any{args.Doc[\"_id\"], args.Doc[\"type\"]})\n\n\treturn out\n}"
    },
    "view-2": {
      "map": "func Map(args couchgo.MapInput) couchgo.MapOutput {  \n    out := couchgo.MapOutput{}\n\tout = append(out, [2]interface{}{args.Doc[\"_id\"], 1})\n\tout = append(out, [2]interface{}{args.Doc[\"_id\"], 2})\n\tout = append(out, [2]interface{}{args.Doc[\"_id\"], 3})\n\t\n\treturn out\n}",
      "reduce": "func Reduce(args couchgo.ReduceInput) couchgo.ReduceOutput {\n\tout := 0.0\n\n\tfor _, value := range args.Values {\n\t\tout += value.(float64)\n\t}\n\n\treturn out\n}"
    }
  },
  "filters": {
    "filter-1": "func Filter(args couchgo.FilterInput) couchgo.FilterOutput {\n\treturn args.Doc[\"type\"] == \"post\"\n}"
  },
  "updates": {
    "update-1": "func Update(args couchgo.UpdateInput) couchgo.UpdateOutput {\n\targs.Doc[\"updated\"] = true\n\treturn couchgo.UpdateOutput{args.Doc, \"ok\"}\n}"
  },
  "rewrites": "func Rewrite(args couchgo.RewriteInput) couchgo.RewriteOutput {\n\treturn couchgo.RewriteOutput{\n\t\tHeaders: map[string]string{\"Location\": \"https://example.com\"},\n\t\tCode:    302,\n\t}\n}",
  "validate_doc_update": "func Validate(args couchgo.ValidateInput) couchgo.ValidateOutput {\n\treturn nil\n}",
  "language": "go"
}
```

## Caveats

The Couchdb Query Server Go is a work in progress and is not yet feature complete.
It implements only minimal functionality without commands that are marked as deprecated in the CouchDB documentation.

## Notes

```shell
go build -o ./bin/couchgo ./lib
```

```shell
# Views
curl 'http://admin:admin@localhost:5984/posts/_design/go/_view/view-1'
curl 'http://admin:admin@localhost:5984/posts/_design/go/_view/view-2'
curl 'http://admin:admin@localhost:5984/posts/_design/go/_view/view-2?reduce=false'

# Filters
curl 'http://admin:admin@localhost:5984/posts/_changes?filter=go/filter-1'
curl 'http://admin:admin@localhost:5984/posts/_changes?filter=_view&view=go/view-1'
curl 'http://admin:admin@localhost:5984/posts/_changes?filter=_view&view=go/view-2'

# Updates
curl -X POST 'http://admin:admin@localhost:5984/posts/_design/go/_update/update-1' -d '{"data": "test"}'
curl -X POST 'http://admin:admin@localhost:5984/posts/_design/go/_update/update-1/70bf54067abb69a75e30ae5f8a0032f1' -d '{"data": "test"}'
```

## Todo

- Improve types, create dedicated types instead of using `interface{}` everywhere
- Improve command handling from stdin so query server can receive properly typed arguments
