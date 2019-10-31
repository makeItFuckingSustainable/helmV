package flags

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// CLIArgs holds all parsed cli flags
type CLIArgs struct {
	Files  []string
	Output string
}

// Parse evaluates the cli flags and parses them into CLIArgs
func Parse() (CLIArgs, error) {
	result := CLIArgs{}
	// TODO include log output
	var f files
	flag.Var(&f, "values", fmt.Sprintf("%s %s",
		"Valid yaml file holding input values.",
		"Multiple value files are processed first to last.",
	))
	flag.Var(&f, "v", fmt.Sprintf("%s %s %s",
		"File holding input values.",
		"Relative path will be changed to absolute path as \"${PWD}/filename\".",
		"Multiple value files are processed first to last.",
	))
	flag.StringVar(&result.Output, "output", "",
		"Absolute output path. Will default to \"${PWD}/values.yaml\".")
	flag.Parse()

	absFiles := make([]string, 0, len(f))
	for _, fileName := range f {
		absFile, err := filepath.Abs(fileName)
		if err != nil {
			return CLIArgs{}, err
		}
		absFiles = append(absFiles, absFile)

	}
	result.Files = absFiles

	err := result.validateOutput()
	if err != nil {
		return CLIArgs{}, err
	}
	err = result.validateFiles()
	if err != nil {
		return CLIArgs{}, err
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

func (args *CLIArgs) validateOutput() error {
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

func (args *CLIArgs) validateFiles() error {
	for _, f := range args.Files {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return fmt.Errorf("value file path \"%s\" does not exist", f)
		}
	}
	return nil
}
