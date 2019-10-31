package process_test

import (
	"bytes"
	"fmt"
	"hval/internal/process"
	"strings"
	"testing"
)

var testValues = []struct {
	deSanitized string
	sanitized   string
	errSan      error
	errDesan    error
}{
	{
		deSanitized: `t1: {{ printf "%s" .k2 }}
t2: {{ .t1 }}{{ .t1 }}
t3: {{ .t1 }}{{ .t1 }} additional text
normal: value
another: normal value`,
		sanitized: `t1: '{{ printf "%s" .k2 }}'
t2: '{{ .t1 }}{{ .t1 }}'
t3: '{{ .t1 }}{{ .t1 }} additional text'
normal: value
another: normal value`,
		errSan:   nil,
		errDesan: nil,
	},
	{
		deSanitized: `{{ .illegal }}: {{ printf "%s" .k2 }}
t2: {{ .t1 }}{{ .t1 }}
t3: {{ .t1 }}{{ .t1 }} additional text
normal: value`,
		sanitized: `{{ .illegal }}: '{{ printf "%s" .k2 }}'
t2: '{{ .t1 }}{{ .t1 }}'
t3: '{{ .t1 }}{{ .t1 }} additional text'
normal: value`,
		errSan: fmt.Errorf("illegal key found in line \"%s\"",
			`{{ .illegal }}: {{ printf "%s" .k2 }}`,
		),
		errDesan: nil,
	},
}

func TestSanitize(t *testing.T) {
	for _, test := range testValues {
		debugSan := new(bytes.Buffer)
		v := process.New(debugSan, true)
		debugDesan := new(bytes.Buffer)
		vOut := process.New(debugDesan, true)
		res, err := v.Sanitize([]byte(test.deSanitized))
		if err != test.errSan {
			if err.Error() != test.errSan.Error() {
				t.Error(errOutput("error", err.Error(), test.errSan.Error()))
			}
		} else {
			if n := strings.Compare(string(res), test.sanitized); n != 0 {
				t.Error(errOutput("sanitization", string(res), test.sanitized))
			}
			resNewline := string(append(res, '\n'))
			if n := strings.Compare(resNewline, debugSan.String()); n != 0 {
				t.Error(errOutput("debugging sanitization",
					debugSan.String(),
					resNewline,
				))
			}
		}

		resDesan, err := vOut.Desanitize([]byte(test.sanitized))
		if err != test.errDesan {
			if err.Error() != test.errDesan.Error() {
				t.Error(errOutput("error", err.Error(), test.errDesan.Error()))
			}
		} else {
			if n := strings.Compare(string(resDesan), test.deSanitized); n != 0 {
				t.Error(errOutput("desanitization", string(resDesan), test.deSanitized))
			}
			resDesanNewline := string(append(resDesan, '\n'))
			if n := strings.Compare(resDesanNewline, debugDesan.String()); n != 0 {
				t.Error(errOutput("debugging desanitization",
					debugDesan.String(),
					resDesanNewline,
				))
			}
		}
	}
}

func errOutput(name, result, expected string) error {
	return fmt.Errorf("[MISMATCH] %s.\nResult: \"%s\" \nExpect: \"%s\"",
		name, result, expected)
}
