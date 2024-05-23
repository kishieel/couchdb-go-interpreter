package main

type Document = map[string]any

type DatabaseInformation struct {
	dbName             string
	committedUpdateSeq int
	docCount           int
	docDelCount        int
	compactRunning     bool
	diskFormatVersion  int
	diskSize           int
	instanceStartTime  string
	purgeSeq           int
	updateSeq          int
}

type SecurityObject struct {
	admins struct {
		names []string
		roles []string
	}
	members struct {
		names []string
		roles []string
	}
}

type UserContext struct {
	db    string
	name  string
	roles []string
}

type Request struct {
	body          string
	cookie        map[string]string
	form          map[string]string
	headers       map[string]string
	id            string
	info          DatabaseInformation
	method        string
	path          []string
	rawPath       string
	requestedPath []string
	peer          string
	query         map[string]string
	secObj        SecurityObject
	userCtx       UserContext
	uuid          string
}

type MapFunction = func(doc Document)
type ReduceFunction = func(keys []any, values []any, rereduce bool) any
type UpdateFunction = func(doc Document, req Request) any
type FilterFunction = func(doc Document, req Request) bool
type ValidateFunction = func(newDoc Document, oldDoc Document, userCtx UserContext, secObj SecurityObject) bool
