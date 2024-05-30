package main

import (
	"fmt"
	"strings"
)

type DesignDocument struct {
	maps     map[string]MapFunction
	reduces  map[string]ReduceFunction
	filters  map[string]FilterFunction
	updates  map[string]UpdateFunction
	rewrite  RewriteFunction
	validate ValidateFunction
}

type QueryServer struct {
	maps    []MapFunction
	reduces []ReduceFunction
	designs map[string]DesignDocument
}

func NewQueryServer() *QueryServer {
	return &QueryServer{
		maps:    []MapFunction{},
		reduces: []ReduceFunction{},
		designs: map[string]DesignDocument{},
	}
}

func (qs *QueryServer) Reset() {
	qs.maps = []MapFunction{}
	qs.reduces = []ReduceFunction{}

	Respond(true)
}

func (qs *QueryServer) AddFun(source string) {
	if strings.HasPrefix(source, "func Map") {
		fn, err := Compile[MapFunction](source)

		if err != nil {
			Respond([]string{"error", "invalid_function", err.Error()})
			return
		}

		qs.maps = append(qs.maps, fn)
	} else if strings.HasPrefix(source, "func Reduce") {
		fn, err := Compile[ReduceFunction](source)

		if err != nil {
			Respond([]string{"error", "invalid_function", err.Error()})
			return
		}

		qs.reduces = append(qs.reduces, fn)
	} else {
		Respond([]string{"error", "invalid_function", "function must be either a map or reduce function"})
		return
	}

	Respond(true)
}

func (qs *QueryServer) AddLib() {
	Log("Libraries are not supported")
	Respond(true)
}

func (qs *QueryServer) MapDoc(doc Document) {
	var results [][][2]any

	for _, fn := range qs.maps {
		results = append(results, fn(MapInput{Doc: doc}))
	}

	Respond(results)
}

func (qs *QueryServer) Reduce(sources []string, keys []any, values []any, rereduce bool) {
	var results []any

	for _, source := range sources {
		fn, err := Compile[ReduceFunction](source)

		if err != nil {
			results = append(results, nil)
			continue
		}

		results = append(results, fn(ReduceInput{keys, values, rereduce}))
	}

	Respond([]any{true, results})
}

func (qs *QueryServer) ProcessDesign(doc map[string]any, path string, compiler func(source string) (any, error), storage map[string]any) {
	if doc[path] != nil {
		for name, source := range doc[path].(map[string]any) {
			fn, err := compiler(source.(string))

			if err != nil {
				Log(fmt.Sprintf("Failed to compile %s function: %s", path, name))
				continue
			}

			storage[name] = fn
		}
	}
}

func (qs *QueryServer) RegisterDesign(docId string, doc map[string]any) {
	maps := make(map[string]MapFunction)
	reduces := make(map[string]ReduceFunction)
	filters := make(map[string]FilterFunction)
	updates := make(map[string]UpdateFunction)
	var rewrite RewriteFunction
	var validate ValidateFunction

	if doc["views"] != nil {
		for name, view := range doc["views"].(map[string]any) {
			if view, ok := view.(map[string]any); ok {
				if source, ok := view["map"].(string); ok {
					fn, err := Compile[MapFunction](source)

					if err != nil {
						Log(fmt.Sprintf("Failed to compile map function: %s", name))
						continue
					}

					maps[name] = fn
				}

				if source, ok := view["reduce"].(string); ok {
					fn, err := Compile[ReduceFunction](source)

					if err != nil {
						Log(fmt.Sprintf("Failed to compile reduce function: %s", name))
						continue
					}

					reduces[name] = fn
				}
			}
		}
	}

	if doc["filters"] != nil {
		for name, source := range doc["filters"].(map[string]any) {
			fn, err := Compile[FilterFunction](source.(string))

			if err != nil {
				Log(fmt.Sprintf("Failed to compile filter function: %s", name))
				continue
			}

			filters[name] = fn
		}
	}

	if doc["updates"] != nil {
		for name, source := range doc["updates"].(map[string]any) {
			fn, err := Compile[UpdateFunction](source.(string))

			if err != nil {
				Log(fmt.Sprintf("Failed to compile update function: %s", name))
				continue
			}

			updates[name] = fn
		}
	}

	if source, ok := doc["rewrites"].(string); ok {
		fn, err := Compile[RewriteFunction](source)

		if err != nil {
			Log("Failed to compile rewrite function")
		}

		rewrite = fn
	}

	if source, ok := doc["validate_doc_update"].(string); ok {
		fn, err := Compile[ValidateFunction](source)

		if err != nil {
			Log("Failed to compile validate function")
		}

		validate = fn
	}

	qs.designs[docId] = DesignDocument{
		maps,
		reduces,
		filters,
		updates,
		rewrite,
		validate,
	}

	Respond(true)
}

func (qs *QueryServer) ExecuteDesign(docId string, path []string, args []any) {
	if path[0] == "filters" {
		fn := qs.designs[docId].filters[path[1]]
		req := Request{} // @todo: fill in request
		results := make([]bool, 0)

		for _, doc := range args[0].([]any) {
			results = append(results, fn(FilterInput{doc.(Document), req}))
		}

		Respond([]any{true, results})
	}

	if path[0] == "views" && path[2] == "map" {
		fn := qs.designs[docId].maps[path[1]]
		results := make([]bool, 0)

		for _, doc := range args[0].([]any) {
			out := fn(MapInput{doc.(Document)})
			results = append(results, len(out) > 0)
		}

		Respond([]any{true, results})
	}

	// rewrite
	// updates
	// validate_doc_update

	Respond([]any{"error", "invalid_path"})
}
