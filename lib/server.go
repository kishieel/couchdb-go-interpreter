package main

import (
	"fmt"
	"strings"
)

type DesignDocument struct {
	views    map[string]MapFunction
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

func (qs *QueryServer) Reset(cmd *ResetCommand) {
	qs.maps = []MapFunction{}
	qs.reduces = []ReduceFunction{}

	Respond(true)
}

func (qs *QueryServer) AddFun(cmd *AddFunCommand) {
	if strings.HasPrefix(cmd.Source, "func Map") {
		fn, err := Compile[MapFunction](cmd.Source)

		if err != nil {
			Respond([]string{"error", "invalid_function", err.Error()})
			return
		}

		qs.maps = append(qs.maps, fn)
	} else if strings.HasPrefix(cmd.Source, "func Reduce") {
		fn, err := Compile[ReduceFunction](cmd.Source)

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

func (qs *QueryServer) AddLib(cmd *AddLibCommand) {
	Log("Libraries are not supported")
	Respond(true)
}

func (qs *QueryServer) MapDoc(cmd *MapDocCommand) {
	var results [][][2]any

	for _, fn := range qs.maps {
		results = append(results, fn(MapInput{cmd.Doc}))
	}

	Respond(results)
}

func (qs *QueryServer) Reduce(cmd *ReduceCommand) {
	var results []any

	for _, source := range cmd.Sources {
		fn, err := Compile[ReduceFunction](source)

		if err != nil {
			results = append(results, nil)
			continue
		}

		results = append(results, fn(ReduceInput{cmd.Keys, cmd.Values, cmd.Rereduce}))
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

func (qs *QueryServer) NewDesign(cmd *NewDesignCommand) {
	views := make(map[string]MapFunction)
	filters := make(map[string]FilterFunction)
	updates := make(map[string]UpdateFunction)
	var rewrite RewriteFunction
	var validate ValidateFunction

	for name, view := range cmd.Doc["views"].(map[string]any) {
		if view, ok := view.(map[string]any); ok {
			if source, ok := view["map"].(string); ok {
				fn, err := Compile[MapFunction](source)

				if err != nil {
					Log(fmt.Sprintf("Failed to compile map function: %s", name))
					continue
				}

				views[name] = fn
			}
		}
	}

	for name, source := range cmd.Doc["filters"].(map[string]any) {
		fn, err := Compile[FilterFunction](source.(string))

		if err != nil {
			Log(fmt.Sprintf("Failed to compile filter function: %s", name))
			continue
		}

		filters[name] = fn
	}

	for name, source := range cmd.Doc["updates"].(map[string]any) {
		fn, err := Compile[UpdateFunction](source.(string))

		if err != nil {
			Log(fmt.Sprintf("Failed to compile update function: %s", name))
			continue
		}

		updates[name] = fn
	}

	if source, ok := cmd.Doc["rewrites"].(string); ok {
		fn, err := Compile[RewriteFunction](source)

		if err != nil {
			Log("Failed to compile rewrite function")
		}

		rewrite = fn
	}

	if source, ok := cmd.Doc["validate_doc_update"].(string); ok {
		fn, err := Compile[ValidateFunction](source)

		if err != nil {
			Log("Failed to compile validate function")
		}

		validate = fn
	}

	qs.designs[cmd.DocId] = DesignDocument{
		views,
		filters,
		updates,
		rewrite,
		validate,
	}

	Respond(true)
}

func (qs *QueryServer) ViewDesign(cmd *ViewDesignCommand) {
	fn := qs.designs[cmd.DocId].views[cmd.FnPth[1]]
	results := make([]bool, 0)

	for _, doc := range cmd.Docs {
		results = append(results, len(fn(MapInput{doc})) > 0)
	}

	Respond([]any{true, results})
}

func (qs *QueryServer) FilterDesign(cmd *FilterDesignCommand) {
	fn := qs.designs[cmd.DocId].filters[cmd.FnPth[1]]
	results := make([]bool, 0)

	for _, doc := range cmd.Docs {
		results = append(results, fn(FilterInput{doc, cmd.Req}))
	}

	Respond([]any{true, results})
}

func (qs *QueryServer) UpdateDesign(cmd *UpdateDesignCommand) {
	fmt.Printf("Command: %v\n", cmd)
}

func (qs *QueryServer) ValidateDesign(cmd *ValidateDesignCommand) {
	fmt.Printf("Command: %v\n", cmd)
}
