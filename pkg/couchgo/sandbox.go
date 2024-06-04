package couchgo

import (
	"github.com/traefik/yaegi/interp"
	"reflect"
)

var Sandbox = interp.Exports{}

func init() {
	Sandbox["couchgo/couchgo"] = map[string]reflect.Value{
		"Log":                 reflect.ValueOf(Log),
		"DatabaseInformation": reflect.ValueOf((*DatabaseInformation)(nil)),
		"Document":            reflect.ValueOf((*Document)(nil)),
		"Request":             reflect.ValueOf((*Request)(nil)),
		"SecurityObject":      reflect.ValueOf((*SecurityObject)(nil)),
		"UserContext":         reflect.ValueOf((*UserContext)(nil)),
		"ForbiddenError":      reflect.ValueOf((*ForbiddenError)(nil)),
		"UnauthorizedError":   reflect.ValueOf((*UnauthorizedError)(nil)),
		"MapInput":            reflect.ValueOf((*MapInput)(nil)),
		"MapOutput":           reflect.ValueOf((*MapOutput)(nil)),
		"ReduceInput":         reflect.ValueOf((*ReduceInput)(nil)),
		"ReduceOutput":        reflect.ValueOf((*ReduceOutput)(nil)),
		"UpdateInput":         reflect.ValueOf((*UpdateInput)(nil)),
		"UpdateOutput":        reflect.ValueOf((*UpdateOutput)(nil)),
		"FilterInput":         reflect.ValueOf((*FilterInput)(nil)),
		"FilterOutput":        reflect.ValueOf((*FilterOutput)(nil)),
		"ValidateInput":       reflect.ValueOf((*ValidateInput)(nil)),
		"ValidateOutput":      reflect.ValueOf((*ValidateOutput)(nil)),
	}
}
