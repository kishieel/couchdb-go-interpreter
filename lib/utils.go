package main

import (
	"encoding/json"
	"fmt"
	"github.com/traefik/yaegi/interp"
	"reflect"
	"strings"
)

var FunctionNames = []string{
	"Map",
	"Reduce",
	"Update",
	"Filter",
	"View",
	"Validate",
	"Rewrite",
}

func Log(message any) {
	if buf, err := json.Marshal(message); err != nil {
		Respond([]string{"log", fmt.Sprintf("Failed to marshal message: %v", err)})
	} else {
		Respond([]string{"log", string(buf)})
	}
}

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

	if err := inter.Use(Sandbox); err != nil {
		return fn, err
	}
	inter.ImportUsed()

	if _, err := inter.Eval(source); err != nil {
		return fn, err
	}

	var val reflect.Value
	var err error

	for _, fnName := range FunctionNames {
		if strings.HasPrefix(source, "func "+fnName) {
			val, err = inter.Eval(fnName)
			break
		}
	}

	if err != nil {
		return fn, err
	}

	fn, ok := val.Interface().(T)
	if !ok {
		return fn, fmt.Errorf("failed to convert function to type")
	}

	return fn, nil
}
