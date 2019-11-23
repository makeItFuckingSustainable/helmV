package main

import (
	"io/ioutil"
	"log"

	"github.com/makeItFuckingSustainable/helmV/internal/cmd/cli"
	"github.com/makeItFuckingSustainable/helmV/internal/cmd/helmV"
	"github.com/makeItFuckingSustainable/helmV/pkg/logerrs"
)

func main() {

	args, err := cli.ParseArgs()
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	e, d := logerrs.New(args.Debug)

	files, err := cli.ReadFiles(args.Files)
	e.Check(err)

	infl, err := helmV.ParseFiles(files, d)
	e.Check(err)
	res, err := helmV.Render(infl, args.MaxIterations, d)
	e.Check(err)

	e.Check(ioutil.WriteFile(args.Output, res, 0666))

}
