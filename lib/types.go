package main

type Document = map[string]any

type DatabaseInformation struct {
	DBName             string `map:"db_name"`
	CommittedUpdateSeq int    `map:"committed_update_seq"`
	DocCount           int    `map:"doc_count"`
	DocDelCount        int    `map:"doc_del_count"`
	CompactRunning     bool   `map:"compact_running"`
	DiskFormatVersion  int    `map:"disk_format_version"`
	DiskSize           int    `map:"disk_size"`
	InstanceStartTime  string `map:"instance_start_time"`
	PurgeSeq           int    `map:"purge_seq"`
	UpdateSeq          int    `map:"update_seq"`
}

type SecurityObject struct {
	Admins struct {
		Names []string `map:"names"`
		Roles []string `map:"roles"`
	} `map:"admins"`
	Members struct {
		Names []string `map:"names"`
		Roles []string `map:"roles"`
	} `map:"members"`
}

type UserContext struct {
	DB    string   `map:"db"`
	Name  string   `map:"name"`
	Roles []string `map:"roles"`
}

type Request struct {
	Body          string              `map:"body"`
	Cookie        map[string]string   `map:"cookie"`
	Form          map[string]string   `map:"form"`
	Headers       map[string]string   `map:"headers"`
	ID            string              `map:"id"`
	Info          DatabaseInformation `map:"info"`
	Method        string              `map:"method"`
	Path          []string            `map:"path"`
	RawPath       string              `map:"raw_path"`
	RequestedPath []string            `map:"requested_path"`
	Peer          string              `map:"peer"`
	Query         map[string]string   `map:"query"`
	SecObj        SecurityObject      `map:"sec_obj"`
	UserCtx       UserContext         `map:"user_ctx"`
	Uuid          string              `map:"uuid"`
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

type FilterOutput = bool

type ValidateInput struct {
	NewDoc  Document
	OldDoc  Document
	UserCtx UserContext
	SecObj  SecurityObject
}

type ValidateOutput = error

type RewriteInput struct {
	Req Request
}

type RewriteOutput struct {
	Path    string
	Query   map[string]string
	Headers map[string]string
	Method  string
	Body    string
	Code    int
}

type MapFunction = func(MapInput) MapOutput
type ReduceFunction = func(ReduceInput) ReduceOutput
type UpdateFunction = func(UpdateInput) UpdateOutput
type FilterFunction = func(FilterInput) FilterOutput
type ValidateFunction = func(ValidateInput) ValidateOutput
type RewriteFunction = func(RewriteInput) RewriteOutput
