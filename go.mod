module github.com/kishieel/couchdb-query-server-go

go 1.21

retract v1.0.0

require (
	github.com/mitchellh/mapstructure v1.5.0
	github.com/traefik/yaegi v0.16.1
)

replace github.com/kishieel/couchgo => ./pkg/couchgo
