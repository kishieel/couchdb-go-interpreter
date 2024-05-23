package main

import (
	"encoding/json"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"reflect"
	"strings"
)

const (
	MapPrefix      = "func Map"
	ReducePrefix   = "func Reduce"
	UpdatePrefix   = "func Update"
	FilterPrefix   = "func Filter"
	ValidatePrefix = "func Validate"
)

var Emitted [][]any

func Respond(message any) {
	if buf, err := json.Marshal(message); err != nil {
		Log(fmt.Sprintf("Error converting object to JSON: %v", err))
		Log(fmt.Sprintf("error on obj: %v", message))
	} else {
		fmt.Println(string(buf))
	}
}

func Compile[T any](source string) (T, error) {
	var fn T

	inter := interp.New(interp.Options{})

	symbols := interp.Exports{
		"couchgo/couchgo": {
			// function, constant and variable definitions
			"Emit": reflect.ValueOf(Emit),
			"Log":  reflect.ValueOf(Log),

			// function, constant and variable definitions
			"DatabaseInformation": reflect.ValueOf((*DatabaseInformation)(nil)),
			"Document":            reflect.ValueOf((*Document)(nil)),
			"Request":             reflect.ValueOf((*Request)(nil)),
			"SecurityObject":      reflect.ValueOf((*SecurityObject)(nil)),
			"UserContext":         reflect.ValueOf((*UserContext)(nil)),
		},
	}

	if err := inter.Use(symbols); err != nil {
		return fn, err
	}

	inter.ImportUsed()

	if _, err := inter.Eval(source); err != nil {
		return fn, err
	}

	var val reflect.Value
	var err error

	switch {
	case strings.HasPrefix(source, MapPrefix):
		val, err = inter.Eval("Map")
	case strings.HasPrefix(source, ReducePrefix):
		val, err = inter.Eval("Reduce")
	case strings.HasPrefix(source, UpdatePrefix):
		val, err = inter.Eval("Update")
	case strings.HasPrefix(source, FilterPrefix):
		val, err = inter.Eval("Filter")
	case strings.HasPrefix(source, ValidatePrefix):
		val, err = inter.Eval("Validate")
	}

	if err != nil {
		return fn, err
	}

	fn, ok := val.Interface().(T)
	if !ok {
		return fn, err
	}

	return fn, nil
}
