package main

import (
	"io/ioutil"
	"log"

	"github.com/makeItFuckingSustainable/helmV/internal/cmd/cli"
	"github.com/makeItFuckingSustainable/helmV/internal/cmd/helmV"
)

func main() {

	args, err := cli.ParseArgs()
	check(err)

	files, err := cli.ReadFiles(args.Files)
	check(err)

	res, err := helmV.Render(files, args.MaxIterations, args.Debug)
	check(err)

	check(ioutil.WriteFile(args.Output, res, 0666))

}

func check(err error) {
	// TODO add proper error and log handling
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
