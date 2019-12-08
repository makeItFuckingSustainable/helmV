package logerrs

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
)

type DebugErr interface {
	Check(error)
	Debugger() Debugger
}

func New(debugFlag bool) DebugErr {
	debugOutput := new(bytes.Buffer)
	d := NewDebugger(debugOutput, debugFlag)
	return debugErr{
		debug:    debugFlag,
		output:   debugOutput,
		debugger: d,
	}
}

type debugErr struct {
	debug    bool
	output   *bytes.Buffer
	debugger Debugger
}

func (e debugErr) Check(err error) {
	// TODO add proper error and log handling
	if err != nil {
		if e.debug {
			fmt.Println("")
			scanner := bufio.NewScanner(e.output)
			for scanner.Scan() {
				log.Printf("DEBUG | %s\n", scanner.Text())
			}
		}
		log.Fatalf("[ERROR] %s", err)
	}
}

func (e debugErr) Debugger() Debugger {
	return debug{
		writer:  e.output,
		doDebug: e.debug,
	}
}
