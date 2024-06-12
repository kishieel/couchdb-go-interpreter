# CouchGO! - CouchDB Query Server

The CouchGO! is a CouchDB Query Server written in Go. It allows you to write your design documents in Go and use them in CouchDB.
It is a drop-in replacement for the JavaScript implementation that ships with CouchDB.

## Usage

To use the CouchGO!, you need to add `couchgo` executable to the CouchDB and add the following
environment variable:

```shell
export COUCHDB_QUERY_SERVER_GO=/path/to/couchgo
```

Now you can create design documents in Go, you just have to set value of `language` property to `go` in the design document.
The structure of the design document is the same as in the JavaScript implementation.

```json
{
  "_id": "_design/ddoc-go",
  "language": "go",
  "views": {
    "example": {
      "map": "func Map(args couchgo.MapInput) couchgo.MapOutput { ... }",
      "reduce": "func Reduce(args couchgo.ReduceInput) couchgo.ReduceOutput { ... }"
    }
  }
}
```

### Design Document

CouchGO! comes with its own approach to design functions. It does not use `emit` or other functions known from the JavaScript implementation, 
but instead expects function to return expected results. Each function accepts only one argument as input and returns a single value (it can be a struct, map, slice, etc.).
Each function can have appropriate name depending on the type of the function.
The details of the functions are described in following sections.

### Map Function

The input of the map function is a single argument of type `MapInput` that contains the document that is being processed.
The output of the map function is a slice of slices of two elements. The first element is the key and the second element is the value.

```go
type MapInput struct{ Doc Document }

type MapOutput [][2]any
```

The following is an example of a map function that emits the document id and type:

```go
func Map(args couchgo.MapInput) couchgo.MapOutput {
    out := make([][2]any, 0)
    out = append(out, [2]any{args.Doc["_id"], args.Doc["type"]})

    return out
}
```

### Reduce Function

The input of the reduce function is a single argument of type `ReduceInput` that contains the keys and values that are being processed and a flag that indicates if the function is being called in the rereduce mode.
The output of the reduce function is a single value that is the result of the reduce operation.

```go
type ReduceInput struct {
	Keys     []any
	Values   []any
	Rereduce bool
}

type ReduceOutput any
```

The following is an example of a reduce function that sums the values:

```go
func Reduce(args couchgo.ReduceInput) couchgo.ReduceOutput {
	out := 0.0
	for _, value := range args.Values {
		out += value.(float64)
	}
	
	return out
}
```

### Filter Function

The input of the filter function is a single argument of type `FilterInput` that contains the document that is being processed and the request object.
The output of the filter function is a boolean value that indicates if the document should be included in the view.

```go
type FilterInput struct {
    Doc Document
    Req Request
}

type FilterOutput = bool
```

The following is an example of a filter function that includes only documents of type `post`:

```go
func Filter(args couchgo.FilterInput) couchgo.FilterOutput {
    return args.Doc["type"] == "post"
}
``` 

### Update Function

The input of the update function is a single argument of type `UpdateInput` that contains the document that is being processed and the request object.
The output of the update function is a single value that is the result of the update operation.

```go
type UpdateInput struct {
    Doc Document
    Req Request
}

type UpdateOutput struct {
    Doc Document
    Res any
}
```

The following is an example of an update function that updates the document:

```go
func Update(args couchgo.UpdateInput) couchgo.UpdateOutput {
    doc := args.Doc
    if doc == nil {
        doc = make(couchgo.Document)
        doc["_id"] = args.Req.Uuid
    }
    
    doc["type"] = "user"
    doc["username"] = "test"
    doc["email"] = "test"
    doc["updated"] = true
    doc["data"] = args.Req.Body

    return couchgo.UpdateOutput{doc, map[string]any{"body": "Updated, ID:" + args.Req.Uuid}}
}
```

### Validate Doc Update Function

The input of the validate doc update function is a single argument of type `ValidateInput` that contains the new document, the old document, the user context, and the security object.
The output of the validate doc update function is an error value that indicates if the document should be updated. Similar to the JavaScript implementation, you can return a `ForbiddenError` and `UnauthorizedError` to indicate that the document should not be updated.

```go
type ValidateInput struct {
    NewDoc Document
    OldDoc Document
    UsrCtx UserContext
    SecObj SecurityObject
}

type ValidateOutput = error
```

The following is an example of a validate doc update function that validates the document type:


```go
func Validate(args couchgo.ValidateInput) couchgo.ValidateOutput {
    if args.NewDoc["type"] == "post" {
        if args.NewDoc["title"] == nil || args.NewDoc["content"] == nil {
            return couchgo.ForbiddenError{Message: "Title and content are required"}
        }

        return nil
    }
    
    if args.NewDoc["type"] == "comment" {
        if args.NewDoc["post"] == nil || args.NewDoc["author"] == nil || args.NewDoc["content"] == nil {
            return couchgo.ForbiddenError{Message: "Post, author, and content are required"}
        }

        return nil
    }
    
    if args.NewDoc["type"] == "user" {
        if args.NewDoc["username"] == nil || args.NewDoc["email"] == nil {
            return couchgo.ForbiddenError{Message: "Username and email are required"}
        }

        return nil
    }
    
    return couchgo.ForbiddenError{Message: "Invalid document type"}
}
```

### Request Object, User Context, and Security Object

The request object, user context, and security object are passed to the filter, update, and validate doc update functions.
The request object contains the following fields:

```go
type Request struct {
	Body          string              
	Cookie        map[string]string   
	Form          map[string]string   
	Headers       map[string]string   
	ID            string              
	Info          DatabaseInformation 
	Method        string              
	Path          []string            
	RawPath       string              
	RequestedPath []string            
	Peer          string              
	Query         map[string]string   
	SecObj        SecurityObject      
	UsrCtx        UserContext         
	Uuid          string              
}
```

The database information contains the following fields:

```go

type DatabaseInformation struct {
	DBName             string `map:"db_name"`
	CommittedUpdateSeq int    `map:"committed_update_seq"`
	DocCount           int    `map:"doc_count"`
	DocDelCount        int    `map:"doc_del_count"`
	CompactRunning     bool   `map:"compact_running"`
	DiskFormatVersion  int    `map:"disk_format_version"`
	DiskSize           int    `map:"disk_size"`
	InstanceStartTime  string `map:"instance_start_time"`
	PurgeSeq           int    `map:"purge_seq"`
	UpdateSeq          int    `map:"update_seq"`
}
```

The user context contains the following fields:

```go
type UserContext struct {
	DB    string   
	Name  string   
	Roles []string 
}
```

The security object contains the following fields:

```go
type SecurityObject struct {
    Admins struct {
        Names []string
        Roles []string
    }
    Members struct {
        Names []string
        Roles []string
    } 
}
```

The document is represented as a map of strings to any value:

```go
type Document map[string]any
```

## Benchmarks

The following benchmarks were performed in a Docker container with 4 CPUs and 2 GB of RAM.
The script used to run the benchmarks is available in the `scripts` directory

| Test                  | CouchGO! |  CouchJS | Boost |
|-----------------------|---------:|---------:|------:|
| Indexing (100k docs)  | 141.713s | 421.529s | 2.97x |
| Reducing (100k docs)  |   7672ms |  15642ms | 2.04x |
| Filtering (100k docs) |  28.928s |  80.594s | 2.79x |
| Updating (1k docs)    |   7.742s |   9.661s | 1.25x |

## Caveats

The CouchGO! is a work in progress and is not yet feature complete.
It implements only minimal functionality without commands that are marked as deprecated in the CouchDB documentation.

CouchGO! does not support the following features yet:

- Show functions
- List functions
- Rewrite functions
