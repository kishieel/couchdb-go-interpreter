package couchgo

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

var FunctionNames = []string{
	"Map",
	"Reduce",
	"Update",
	"Filter",
	"View",
	"Validate",
}
