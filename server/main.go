package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	server := NewQueryServer()

	dispatch := map[string]func(args ...any){
		"ddoc":      func(args ...any) {},
		"reset":     func(args ...any) { server.Reset() },
		"add_fun":   func(args ...any) { server.AddFun(args[0].(string)) },
		"add_lib":   func(args ...any) {},
		"map_doc":   func(args ...any) {},
		"index_doc": func(args ...any) {},
		"reduce":    func(args ...any) {},
		"rereduce":  func(args ...any) {},
		"test":      func(args ...any) { Log("Welcome from Go!") },
	}

	for scanner.Scan() {
		var command []any

		if err := json.Unmarshal(scanner.Bytes(), &command); err != nil {
			Respond([]string{"error", "unnamed_error", err.Error()})
		}

		cmdkey := command[0].(string)
		if dispatch[cmdkey] != nil {
			dispatch[cmdkey](command[1:]...)
		} else {
			Respond([]string{"error", "unknown_command", fmt.Sprintf("unknown command '%s'", cmdkey)})
		}
	}
}
