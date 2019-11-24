package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ReadFiles iterates through all filePaths and parses the files to byte slices.
func ReadFiles(filePaths []string) ([][]byte, error) {
	res := make([][]byte, 0, len(filePaths))
	for _, path := range filePaths {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return [][]byte{},
				fmt.Errorf("cannot read file \"%s\" - error: %s", f, err)
		}
		res = append(res, f)
	}
	return res, nil
}

// Args holds all parsed cli flags
type Args struct {
	Files         []string
	Output        string
	MaxIterations uint
	Debug         bool
}

// ParseArgs evaluates the cli flags and parses them into CLIArgs
func ParseArgs() (Args, error) {
	result := Args{}
	// TODO include log output
	var f files
	flag.Var(
		&f,
		"values",
		fmt.Sprintf("%s %s",
			"File holding input values.",
			"Multiple value files are processed first to last.",
		))
	flag.Var(
		&f,
		"v",
		fmt.Sprintf("%s %s %s",
			"File holding input values.",
			"Relative path will be changed to absolute path as \"${PWD}/filename\".",
			"Multiple value files are processed first to last.",
		))
	flag.StringVar(
		&result.Output,
		"output",
		"output.yaml",
		"Absolute output path. Will default to \"${PWD}/values.yaml\".",
	)
	flag.UintVar(
		&result.MaxIterations,
		"max-iterations",
		10,
		"Maximal number of recursive iterations that helmV will execute.",
	)
	flag.BoolVar(
		&result.Debug,
		"debug",
		false,
		"Activate debugging.",
	)
	flag.Parse()

	absFiles := make([]string, 0, len(f))
	for _, fileName := range f {
		absFile, err := filepath.Abs(fileName)
		if err != nil {
			return Args{}, err
		}
		absFiles = append(absFiles, absFile)

	}
	result.Files = absFiles

	err := result.validateOutput()
	if err != nil {
		return Args{}, err
	}
	err = result.validateFiles()
	if err != nil {
		return Args{}, err
	}
	return result, nil
}

type files []string

func (f *files) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *files) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (args *Args) validateOutput() error {
	if args.Output == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		args.Output = fmt.Sprintf("%s/values.yaml", cwd)
	}
	if _, err := os.Stat(args.Output); err == nil {
		return fmt.Errorf("output file path \"%s\" does already exist", args.Output)
	}
	return nil
}

func (args *Args) validateFiles() error {
	for _, f := range args.Files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return fmt.Errorf("value file path \"%s\" does not exist", f)
		}
	}
	return nil
}
