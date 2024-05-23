package main

import (
	"strings"
)

type QueryServer struct {
	functions []MapFunction
}

func NewQueryServer() *QueryServer {
	return &QueryServer{
		functions: []MapFunction{},
	}
}

func (qs *QueryServer) Reset() {
	qs.functions = []MapFunction{}
	Respond(true)
}

func (qs *QueryServer) AddFun(source string) {
	var fn interface{}
	var err error

	if strings.HasPrefix(source, "func Map") {
		fn, err = Compile[MapFunction](source)
	} else if strings.HasPrefix(source, "func Reduce") {
		fn, err = Compile[ReduceFunction](source)
	} else {
		Respond([]string{"error", "invalid_function_type", "Invalid function type"})
		return
	}

	if err != nil {
		Respond([]string{"error", "compilation_error", err.Error()})
		return
	}

	if fn, ok := fn.(MapFunction); ok {
		qs.functions = append(qs.functions, fn)
	}

	Respond(true)
}

func (qs *QueryServer) AddLib() {
	Respond(true)
}

func (qs *QueryServer) MapDoc(doc Document) {
	var results [][][]any

	for _, fn := range qs.functions {
		Emitted = [][]any{}
		fn(doc)
		results = append(results, Emitted)
	}

	Respond(results)
}

func (qs *QueryServer) Reduce(sources []string, keys []any, values []any, rereduce bool) {
	var reductions []any

	for _, source := range sources {
		fn, err := Compile[ReduceFunction](source)

		if err != nil {
			reductions = append(reductions, nil)
			continue
		}

		reductions = append(reductions, fn(keys, values, rereduce))
	}

	Respond([]any{true, reductions})
}
