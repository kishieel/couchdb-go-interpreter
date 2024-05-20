package main

import "runtime"

type QueryServer struct {
	functions []func(doc map[string]any)
}

func NewQueryServer() *QueryServer {
	return &QueryServer{}
}

func (qs *QueryServer) Reset() {
	qs.functions = []func(doc map[string]any){}
	runtime.GC()
	Respond("true")
}

func (qs *QueryServer) AddFun(source string) {
	symbol, err := CompileFunction(source)

	if err != nil {
		Respond([]string{"error", "compilation_error", err.Error()})
		return
	}

	fn, ok := symbol.(func(doc map[string]any))

	if !ok {
		Respond([]string{"error", "type_assertion_failed", "symbol is not a function"})
		return
	}

	qs.functions = append(qs.functions, fn)
	Respond("true")
}
