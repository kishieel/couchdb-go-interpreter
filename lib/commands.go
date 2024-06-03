package main

import "github.com/mitchellh/mapstructure"

type CommandKind string

const (
	Reset          CommandKind = "reset"
	AddFun         CommandKind = "add_fun"
	AddLib         CommandKind = "add_lib"
	MapDoc         CommandKind = "map_doc"
	Reduce         CommandKind = "reduce"
	Rereduce       CommandKind = "rereduce"
	Design         CommandKind = "ddoc"
	NewDesign      CommandKind = "new_design"
	ViewDesign     CommandKind = "view_design"
	FilterDesign   CommandKind = "filter_design"
	UpdateDesign   CommandKind = "update_design"
	ValidateDesign CommandKind = "validate_design"
)

type Command interface {
	Parse(args ...any)
}

type CommandFactory func(args ...any) Command

var CommandRegistry = map[CommandKind]CommandFactory{
	Reset:          func(args ...any) Command { return &ResetCommand{} },
	AddFun:         func(args ...any) Command { return &AddFunCommand{} },
	AddLib:         func(args ...any) Command { return &AddLibCommand{} },
	MapDoc:         func(args ...any) Command { return &MapDocCommand{} },
	Reduce:         func(args ...any) Command { return &ReduceCommand{} },
	Rereduce:       func(args ...any) Command { return &ReduceCommand{} },
	NewDesign:      func(args ...any) Command { return &NewDesignCommand{} },
	ViewDesign:     func(args ...any) Command { return &ViewDesignCommand{} },
	FilterDesign:   func(args ...any) Command { return &FilterDesignCommand{} },
	UpdateDesign:   func(args ...any) Command { return &UpdateDesignCommand{} },
	ValidateDesign: func(args ...any) Command { return &ValidateDesignCommand{} },
}

var DesignOperationRegistry = map[string]CommandKind{
	"views":               ViewDesign,
	"filters":             FilterDesign,
	"updates":             UpdateDesign,
	"validate_doc_update": ValidateDesign,
}

var CommandDispatcher = map[CommandKind]func(*QueryServer, Command){
	Reset:          func(server *QueryServer, cmd Command) { server.Reset(cmd.(*ResetCommand)) },
	AddFun:         func(server *QueryServer, cmd Command) { server.AddFun(cmd.(*AddFunCommand)) },
	AddLib:         func(server *QueryServer, cmd Command) { server.AddLib(cmd.(*AddLibCommand)) },
	MapDoc:         func(server *QueryServer, cmd Command) { server.MapDoc(cmd.(*MapDocCommand)) },
	Reduce:         func(server *QueryServer, cmd Command) { server.Reduce(cmd.(*ReduceCommand)) },
	Rereduce:       func(server *QueryServer, cmd Command) { server.Reduce(cmd.(*ReduceCommand)) },
	NewDesign:      func(server *QueryServer, cmd Command) { server.NewDesign(cmd.(*NewDesignCommand)) },
	ViewDesign:     func(server *QueryServer, cmd Command) { server.ViewDesign(cmd.(*ViewDesignCommand)) },
	FilterDesign:   func(server *QueryServer, cmd Command) { server.FilterDesign(cmd.(*FilterDesignCommand)) },
	UpdateDesign:   func(server *QueryServer, cmd Command) { server.UpdateDesign(cmd.(*UpdateDesignCommand)) },
	ValidateDesign: func(server *QueryServer, cmd Command) { server.ValidateDesign(cmd.(*ValidateDesignCommand)) },
}

type ResetCommand struct {
	Kind   CommandKind
	Config map[string]interface{}
}

func (cmd *ResetCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.Config = args[1].(map[string]interface{})
}

type AddFunCommand struct {
	Kind   CommandKind
	Source string
}

func (cmd *AddFunCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.Source = args[1].(string)
}

type AddLibCommand struct {
	Kind CommandKind
	Lib  map[string]interface{}
}

func (cmd *AddLibCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.Lib = args[1].(map[string]interface{})
}

type MapDocCommand struct {
	Kind CommandKind
	Doc  map[string]interface{}
}

func (cmd *MapDocCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.Doc = args[1].(map[string]interface{})
}

type ReduceCommand struct {
	Kind     CommandKind
	Sources  []string
	Keys     []interface{}
	Values   []interface{}
	Rereduce bool
}

func (cmd *ReduceCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.Sources = make([]string, len(args[1].([]any)))
	cmd.Keys = make([]interface{}, len(args[2].([]any)))
	cmd.Values = make([]interface{}, len(args[2].([]any)))

	for i, source := range args[1].([]any) {
		cmd.Sources[i] = source.(string)
	}

	if cmd.Kind == Rereduce {
		cmd.Rereduce = true
		cmd.Keys = nil
		cmd.Values = args[2].([]any)
	} else {
		cmd.Rereduce = false
		for i, kv := range args[2].([]any) {
			cmd.Keys[i] = kv.([]any)[0]
			cmd.Values[i] = kv.([]any)[1]
		}
	}
}

type NewDesignCommand struct {
	Kind  CommandKind
	DocId string
	Doc   map[string]interface{}
}

func (cmd *NewDesignCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.DocId = args[2].(string)
	cmd.Doc = args[3].(map[string]interface{})
}

type ViewDesignCommand struct {
	Kind  CommandKind
	DocId string
	FnPth [3]string
	Docs  []map[string]interface{}
}

func (cmd *ViewDesignCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.DocId = args[1].(string)
	cmd.FnPth[0] = args[2].([]any)[0].(string)
	cmd.FnPth[1] = args[2].([]any)[1].(string)
	cmd.FnPth[2] = args[2].([]any)[2].(string)

	cmd.Docs = make([]map[string]interface{}, len(args[3].([]any)))
	for i, doc := range args[3].([]any) {
		cmd.Docs[i] = doc.(map[string]interface{})
	}
}

type FilterDesignCommand struct {
	Kind  CommandKind
	DocId string
	FnPth [2]string
	Docs  []map[string]interface{}
	Req   Request
}

func (cmd *FilterDesignCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.DocId = args[1].(string)
	cmd.FnPth[0] = args[2].([]any)[0].(string)
	cmd.FnPth[1] = args[2].([]any)[1].(string)
	cmd.Docs = make([]map[string]interface{}, len(args[3].([]any)))
	for i, doc := range args[3].([]any) {
		cmd.Docs[i] = doc.(map[string]interface{})
	}
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &cmd.Req})
	decoder.Decode(args[4])
}

type UpdateDesignCommand struct {
	Kind  CommandKind
	DocId string
	FnPth [2]string
	Doc   map[string]interface{}
	Req   Request
}

func (cmd *UpdateDesignCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.DocId = args[1].(string)
	cmd.FnPth[0] = args[2].([]any)[0].(string)
	cmd.FnPth[1] = args[2].([]any)[1].(string)
	cmd.Doc = args[3].([]any)[0].(map[string]interface{})

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &cmd.Req})
	decoder.Decode(args[3].([]any)[1])
}

type ValidateDesignCommand struct {
	Kind   CommandKind
	DocId  string
	NewDoc map[string]interface{}
	OldDoc map[string]interface{}
	UsrCtx UserContext
	SecObj SecurityObject
}

func (cmd *ValidateDesignCommand) Parse(args ...any) {
	cmd.Kind = CommandKind(args[0].(string))
	cmd.DocId = args[1].(string)
	cmd.NewDoc = args[3].([]any)[0].(map[string]interface{})
	cmd.OldDoc = args[3].([]any)[1].(map[string]interface{})

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "map", Result: &cmd.UsrCtx})
	decoder.Decode(args[3].([]any)[2])

	decoder, _ = mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "map", Result: &cmd.SecObj})
	decoder.Decode(args[3].([]any)[3])
}

func GetCommandKind(args ...any) CommandKind {
	kind := CommandKind(args[0].(string))

	if kind != Design {
		return kind
	}

	if args[1].(string) == "new" {
		return NewDesign
	}

	return DesignOperationRegistry[args[2].([]any)[0].(string)]
}
