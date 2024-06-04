package main

import (
	"github.com/kishieel/couchdb-query-server-go/pkg/couchgo"
)

func main() {
	server := couchgo.NewQueryServer()
	server.Start()
}
