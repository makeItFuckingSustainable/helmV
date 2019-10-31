package process

import "io"

type input struct {
	debugger io.Writer
	debug    bool
}

func New(debugger io.Writer, debugFlag bool) input {
	return input{
		debugger: debugger,
		debug:    debugFlag,
	}
}
