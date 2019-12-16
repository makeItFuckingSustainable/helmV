package yamltmpl

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
)

var keyIsTmpl = regexp.MustCompile(`({{.*}}.*:.*)`)
var valIsTmpl = regexp.MustCompile(`(\w*: )({{.*}}.*)`)
var sanTmpl = regexp.MustCompile("(\\w*: )'({{.*}}.*)'")

// Sanitize takes a yaml-golang template in byte slice format transforms and
// returns a sanitized yaml-file version of it also in byte slice format.
// All values that are golang-templates in the input are transformed to strings
// in the output.
func Sanitize(input []byte) ([]byte, error) {
	// TODO: transform input to *Scanner type

	res := make([]byte, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(input))
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(keyIsTmpl.FindAll(line, -1)) > 0 {
			return res, fmt.Errorf("illegal key found in line \"%s\"", line)
		}
		// stringify all occurances of template inputs
		matchTmpl := valIsTmpl.ReplaceAll(line, []byte("${1}'${2}'"))
		matchTmpl = append(matchTmpl, '\n')
		res = append(res, matchTmpl...)
	}
	if err := scanner.Err(); err != nil {
		return res, err
	}
	return res[:len(res)-1], nil
}

// Desanitize reverts the action of Sanitize. It takes a yaml-file in byte slice
// format and returns a yaml-golang template version of it also in byte slice format.
// All values that are stringified golang-templates in the input are transformed
// to actual golang-template values in the output.
func Desanitize(input []byte) ([]byte, error) {

	res := make([]byte, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(input))
	for scanner.Scan() {
		line := scanner.Bytes()
		// revert stringification of all template inputs
		matchTmpl := sanTmpl.ReplaceAll(line, []byte("${1}${2}"))
		matchTmpl = append(matchTmpl, '\n')
		res = append(res, matchTmpl...)
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}
	return res[:len(res)-1], nil
}
