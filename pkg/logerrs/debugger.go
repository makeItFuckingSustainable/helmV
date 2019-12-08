package logerrs

import "io"

// Debugger is an interface that holds a debug mode flag (DoDebug) and a debugging
// target (Writer) where the debug information is sent to
type Debugger interface {
	Write([]byte) error
	DoDebug() bool
	SetDebugMode(bool)
}

// NewDebugger constructs a Debugger from a Writer and a flag indicating whether
// debug mode is enabled
func NewDebugger(debugger io.Writer, debugFlag bool) Debugger {
	d := &debug{
		writer:  debugger,
		doDebug: debugFlag,
	}
	return d
}

type debug struct {
	writer  io.Writer
	doDebug bool
}

// Write implements the Write method on the Debugger Writer
func (d debug) Write(dump []byte) error {
	if _, err := d.writer.Write(dump); err != nil {
		return err
	}
	return nil
}

// DoDebug returns a bool indicating whether debug mode is enabled
func (d debug) DoDebug() bool {
	if d.doDebug {
		return true
	}
	return false
}

// SetDebugMode allows to change the debug mode (true or false)
func (d *debug) SetDebugMode(mode bool) {
	d.doDebug = mode
}
