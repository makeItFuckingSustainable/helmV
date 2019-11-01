package main

import (
	"bytes"
	"helmV/internal/cmd"
	"helmV/internal/debug"
	"helmV/pkg/flags"
	"log"
	"os"
)

func main() {

	args, err := flags.Parse()
	check(err)
	d := debug.New(new(bytes.Buffer), false)
	infl, err := cmd.LoadInput(args.Files, d)
	check(err)
	res, err := os.Create("output.yaml")
	check(err)

	check(cmd.RenderResult(infl, res, 10, d))

}

func check(err error) {
	// TODO add proper error and log handling
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
