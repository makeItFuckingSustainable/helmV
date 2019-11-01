package debug

import "io"

type Debugger interface {
	Write([]byte) error
	DoDebug() bool
}

func New(debugger io.Writer, debugFlag bool) debug {
	return debug{
		writer:  debugger,
		doDebug: debugFlag,
	}
}

type debug struct {
	writer  io.Writer
	doDebug bool
}

func (d debug) Write(dump []byte) error {
	if _, err := d.writer.Write(dump); err != nil {
		return err
	}
	return nil
}

func (d debug) DoDebug() bool {
	if d.doDebug {
		return true
	}
	return false
}
