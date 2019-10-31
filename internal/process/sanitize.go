package process

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
)

var keyIsTmpl = regexp.MustCompile(`({{.*}}.*:.*)`)
var valIsTmpl = regexp.MustCompile(`(\w*: )({{.*}}.*)`)
var sanTmpl = regexp.MustCompile("(\\w*: )'({{.*}}.*)'")

func (i *input) Sanitize(input []byte) ([]byte, error) {
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
		if i.debug {
			if _, err := i.debugger.Write(matchTmpl); err != nil {
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

func (i *input) Desanitize(input []byte) ([]byte, error) {
	res := make([]byte, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(input))
	for scanner.Scan() {
		line := scanner.Bytes()
		// revert stringification of all template inputs
		matchTmpl := sanTmpl.ReplaceAll(line, []byte("${1}${2}"))
		matchTmpl = append(matchTmpl, '\n')
		if i.debug {
			if _, err := i.debugger.Write(matchTmpl); err != nil {
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
