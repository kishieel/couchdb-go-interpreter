package main

import (
	"encoding/json"
	"fmt"
)

func Emit(key any, value any) {
	Emitted = append(Emitted, []any{key, value})
}

func Log(message any) {
	if buf, err := json.Marshal(message); err != nil {
		Respond([]string{"log", fmt.Sprintf("Failed to marshal message: %v", err)})
	} else {
		Respond([]string{"log", string(buf)})
	}
}
