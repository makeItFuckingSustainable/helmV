package yamltmpl_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/makeItFuckingSustainable/helmV/internal/yamltmpl"
)

var testValues = []struct {
	deSanitized string
	sanitized   string
	errSan      error
	errDesan    error
	bijective   bool
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
		errSan:    nil,
		errDesan:  nil,
		bijective: true,
	},
	{
		deSanitized: `t2: {{ .t1 }}{{ .t1 }}
{{ .illegal }}: {{ printf "%s" .k2 }}
t3: {{ .t1 }}{{ .t1 }} additional text
normal: value`,
		sanitized: `t2: '{{ .t1 }}{{ .t1 }}'
`,
		errSan: fmt.Errorf("illegal key found in line \"%s\"",
			`{{ .illegal }}: {{ printf "%s" .k2 }}`,
		),
		errDesan:  nil,
		bijective: false,
	},
}

func TestSanitize(t *testing.T) {
	for _, test := range testValues {
		res, err := yamltmpl.Sanitize([]byte(test.deSanitized))
		if err != test.errSan {
			if err.Error() != test.errSan.Error() {
				t.Error(errOutput("error", err.Error(), test.errSan.Error()))
			}
		}

		if n := strings.Compare(string(res), test.sanitized); n != 0 {
			t.Error(errOutput("sanitization", string(res), test.sanitized))
		}

		resDesan, err := yamltmpl.Desanitize([]byte(test.sanitized))
		if err != test.errDesan {
			if err.Error() != test.errDesan.Error() {
				t.Error(errOutput("error", err.Error(), test.errDesan.Error()))
			}
		}
		if test.bijective {
			if n := strings.Compare(string(resDesan), test.deSanitized); n != 0 {
				t.Error(errOutput("desanitization", string(resDesan), test.deSanitized))
			}
		}

	}
}

func errOutput(name, result, expected string) error {
	return fmt.Errorf("[MISMATCH] %s.\nResult: \"%s\" \nExpect: \"%s\"",
		name, result, expected)
}
