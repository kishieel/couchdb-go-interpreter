package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"plugin"
	"strings"
)

var logger = slog.Default()

func Log(message any) {
	if buf, err := json.Marshal(message); err != nil {
		logger.Error("Failed to marshal message", err)
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

func CompileFunction(source string) (plugin.Symbol, error) {
	file, err := os.CreateTemp("", "couchdb-*.go")

	if err != nil {
		logger.Error("Failed to create temporary file", err)
		return nil, err
	}

	source = fmt.Sprintf("package main\n\n%s", source)
	if _, err := file.WriteString(source); err != nil {
		logger.Error("Failed to write source to temporary file", err)
		return nil, err
	}

	if err := file.Close(); err != nil {
		logger.Error("Failed to close temporary file", err)
		return nil, err
	}

	if _, err := exec.Command("go", "build", "-buildmode=plugin", "-o", strings.Replace(file.Name(), ".go", ".so", 1), file.Name()).Output(); err != nil {
		logger.Error("Failed to build plugin", err)
		return nil, err
	}

	plug, err := plugin.Open(strings.Replace(file.Name(), ".go", ".so", 1))
	if err != nil {
		logger.Error("Failed to open plugin", err)
		return nil, err
	}

	return plug.Lookup("Handle")
}
