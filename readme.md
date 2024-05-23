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
  "language": "go",
  "views": {
    "example": {
      "map": "func Handle(doc *Document) { Emit((*doc)[\"_id\"], nil) }",
      "reduce": "func Handle(keys []string, values []interface{}, rereduce bool) interface{} { return len(keys) }"
    }
  }
}
```

## Caveats

The Couchdb Query Server Go is a work in progress and is not yet feature complete.
It implements only minimal functionality without commands that are marked as deprecated in the CouchDB documentation.

## Notes

```shell
go build -o ./bin/couchgo ./lib
```

```json
{
  "_id": "_design/go",
  "views": {
    "by-author": {
      "map": "func Map(doc *couchgo.Document) {\n    if (*doc)[\"type\"] == \"post\" && (*doc)[\"author\"] != nil {\n\t\tcouchgo.Emit((*doc)[\"author\"], nil)\n\t}\n}"
    }
  },
  "language": "go"
}
```

## Todo

- Better error handling for functions executed in sandbox