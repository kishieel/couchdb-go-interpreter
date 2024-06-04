package couchgo

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

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

func (qs *QueryServer) Start() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var event []interface{}

		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			Respond([]string{"error", "unnamed_error", err.Error()})
		}

		Log(event)
		kind := GetCommandKind(event...)

		factory, found := CommandRegistry[kind]
		if !found {
			Respond([]string{"error", "unknown_command", fmt.Sprintf("Unknown command type: %s", kind)})
			continue
		}

		command := factory()
		command.Create(event...)

		dispatcher, found := CommandDispatcher[kind]
		if !found {
			Respond([]string{"error", "unknown_command", fmt.Sprintf("Unknown command type: %s", kind)})
			continue
		}

		dispatcher(qs, command)
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
	var validate ValidateFunction

	if cmd.Doc["views"] != nil {
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
	}

	if cmd.Doc["filters"] != nil {
		for name, source := range cmd.Doc["filters"].(map[string]any) {
			fn, err := Compile[FilterFunction](source.(string))

			if err != nil {
				Log(fmt.Sprintf("Failed to compile filter function: %s", name))
				continue
			}

			filters[name] = fn
		}
	}

	if cmd.Doc["updates"] != nil {
		for name, source := range cmd.Doc["updates"].(map[string]any) {
			fn, err := Compile[UpdateFunction](source.(string))

			if err != nil {
				Log(fmt.Sprintf("Failed to compile update function: %s", name))
				continue
			}

			updates[name] = fn
		}
	}

	if cmd.Doc["validate_doc_update"] != nil {
		if source, ok := cmd.Doc["validate_doc_update"].(string); ok {
			fn, err := Compile[ValidateFunction](source)

			if err != nil {
				Log("Failed to compile validate function")
			}

			validate = fn
		}
	}

	qs.designs[cmd.DocId] = DesignDocument{
		views,
		filters,
		updates,
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
	fn := qs.designs[cmd.DocId].updates[cmd.FnPth[1]]
	result := fn(UpdateInput{cmd.Doc, cmd.Req})

	Respond([]any{"up", result.Doc, result.Res})
}

func (qs *QueryServer) ValidateDesign(cmd *ValidateDesignCommand) {
	fn := qs.designs[cmd.DocId].validate
	err := fn(ValidateInput{cmd.NewDoc, cmd.OldDoc, cmd.UsrCtx, cmd.SecObj})

	if errors.Is(err, UnauthorizedError{}) {
		Respond([]any{"error", "unauthorized", err.Error()})
		return
	}

	if errors.Is(err, ForbiddenError{}) {
		Respond([]any{"error", "forbidden", err.Error()})
		return
	}

	Respond(1)
}
