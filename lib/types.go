package main

type Document = map[string]any

type DatabaseInformation struct {
	DBName             string
	CommittedUpdateSeq int
	DocCount           int
	DocDelCount        int
	CompactRunning     bool
	DiskFormatVersion  int
	DiskSize           int
	InstanceStartTime  string
	PurgeSeq           int
	UpdateSeq          int
}

type SecurityObject struct {
	Admins struct {
		Names []string
		Roles []string
	}
	Members struct {
		Names []string
		Roles []string
	}
}

type UserContext struct {
	DB    string
	Name  string
	Roles []string
}

type Request struct {
	Body          string
	Cookie        map[string]string
	Form          map[string]string
	Headers       map[string]string
	ID            string
	Info          DatabaseInformation
	Method        string
	Path          []string
	RawPath       string
	RequestedPath []string
	Peer          string
	Query         map[string]string
	SecObj        SecurityObject
	UserCtx       UserContext
	Uuid          string
}

type MapInput struct{ Doc Document }

type MapOutput [][2]any

type ReduceInput struct {
	Keys     []any
	Values   []any
	Rereduce bool
}

type ReduceOutput any

type UpdateInput struct {
	Doc Document
	Req Request
}

type UpdateOutput struct {
	Doc Document
	Res any
}

type FilterInput struct {
	Doc Document
	Req Request
}

type FilterOutput bool

type ViewInput struct {
	Doc Document
	Req Request
}

type ViewOutput [][2]any

type ValidateInput struct {
	NewDoc  Document
	OldDoc  Document
	UserCtx UserContext
	SecObj  SecurityObject
}

type ValidateOutput error

type RewriteInput struct {
	Req Request
}

type RewriteOutput struct {
	Path    string
	Query   []string
	Headers map[string]string
	Method  string
	Body    string
	Code    int
}

type MapFunction = func(MapInput) MapOutput
type ReduceFunction = func(ReduceInput) ReduceOutput
type UpdateFunction = func(UpdateInput) UpdateOutput
type FilterFunction = func(FilterInput) FilterOutput
type ViewFunction = func(ViewInput) ViewOutput
type ValidateFunction = func(ValidateInput) ValidateOutput
type RewriteFunction = func(RewriteInput) RewriteOutput
