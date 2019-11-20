package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/makeItFuckingSustainable/helmV/internal/cmd"
	"github.com/makeItFuckingSustainable/helmV/internal/debug"
	"github.com/makeItFuckingSustainable/helmV/pkg/flags"
	"log"
	"os"
)

func main() {

	debugOutput := new(bytes.Buffer)
	e := err{true, debugOutput}
	args, err := flags.Parse()
	e.check(err)
	d := debug.New(debugOutput, args.Debug)
	e.debug = d.DoDebug()

	infl, err := cmd.LoadInput(args.Files, d)
	e.check(err)
	res, err := os.Create(args.Output)
	e.check(err)

	e.check(cmd.RenderResult(infl, res, args.MaxIterations, d))

}

type err struct {
	debug  bool
	output *bytes.Buffer
}

func (e err) check(err error) {
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
