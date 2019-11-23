package logerrs

import (
	"bufio"
	"bytes"
	"fmt"
	"log"

)

type DebugErr interface {
	Check(error)
}

func New(debugFlag bool) (DebugErr, Debugger) {
	debugOutput := new(bytes.Buffer)
	d := NewDebugger(debugOutput, debugFlag)
	return debugErr{
		debug:  debugFlag,
		output: debugOutput,
	}, d
}

type debugErr struct {
	debug  bool
	output *bytes.Buffer
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
