package main

import (
	"strings"
)

type DesignDocument struct {
	filters  map[string]FilterFunction
	views    map[string]ViewFunction
	updates  map[string]UpdateFunction
	rewrites map[string]RewriteFunction
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
