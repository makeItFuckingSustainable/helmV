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
	e = logerrs.New(args.Debug)

	files, err := cli.ReadFiles(args.Files)
	e.Check(err)

	infl, err := helmV.ParseFiles(files, e.Debugger())
	e.Check(err)
	res, err := helmV.Render(infl, args.MaxIterations, e.Debugger())
	e.Check(err)

	e.Check(ioutil.WriteFile(args.Output, res, 0666))

}
