package yamltmpl

import (
	"bufio"
	"bytes"
	"fmt"
	"helmV/internal/debug"
	"regexp"
)

var keyIsTmpl = regexp.MustCompile(`({{.*}}.*:.*)`)
var valIsTmpl = regexp.MustCompile(`(\w*: )({{.*}}.*)`)
var sanTmpl = regexp.MustCompile("(\\w*: )'({{.*}}.*)'")

func Sanitize(input []byte, debug debug.Debugger) ([]byte, error) {
	res := make([]byte, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(input))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(keyIsTmpl.FindAll(line, -1)) > 0 {
			return []byte{}, fmt.Errorf(
				"illegal key found in line \"%s\"", line,
			)
		}
		// stringify all occurances of template inputs
		matchTmpl := valIsTmpl.ReplaceAll(line, []byte("${1}'${2}'"))
		matchTmpl = append(matchTmpl, '\n')
		if debug.DoDebug() {
			if err := debug.Write(matchTmpl); err != nil {
				return []byte{}, err
			}
		}
		res = append(res, matchTmpl...)
	}
	if err := scanner.Err(); err != nil {
		return []byte{}, err
	}
	return res[:len(res)-1], nil
}

func Desanitize(input []byte, debug debug.Debugger) ([]byte, error) {
	res := make([]byte, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(input))
	for scanner.Scan() {
		line := scanner.Bytes()
		// revert stringification of all template inputs
		matchTmpl := sanTmpl.ReplaceAll(line, []byte("${1}${2}"))
		matchTmpl = append(matchTmpl, '\n')
		if debug.DoDebug() {
			if err := debug.Write(matchTmpl); err != nil {
				return []byte{}, err
			}
		}
		res = append(res, matchTmpl...)
	}

	if err := scanner.Err(); err != nil {
		return []byte{}, err
	}
	return res[:len(res)-1], nil
}

func errTmplFormat(line string) error {
	return fmt.Errorf("unexpected template format in line \"%s\"", line)
}
