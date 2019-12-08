package main

import (
	"io/ioutil"

	"github.com/makeItFuckingSustainable/helmV/internal/cmd/cli"
	"github.com/makeItFuckingSustainable/helmV/internal/cmd/helmV"
	"github.com/makeItFuckingSustainable/helmV/pkg/logerrs"
)

func main() {

	e := logerrs.New(true)
	args, err := cli.ParseArgs()
	e.Check(err)
	e.SetDebugMode(args.Debug)

	files, err := cli.ReadFiles(args.Files)
	e.Check(err)

	hV := helmV.New(e.Debugger())
	res, err := hV.Render(files, args.MaxIterations)
	e.Check(err)

	e.Check(ioutil.WriteFile(args.Output, res, 0666))

}
