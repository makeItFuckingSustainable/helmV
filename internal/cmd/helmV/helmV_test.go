package helmV_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/makeItFuckingSustainable/helmV/internal/cmd/helmV"
)

var testValues = []struct {
	name          string
	input         [][]byte
	maxIterations uint
	res           []byte
	err           error
}{
	{
		name: "illegal key template",
		input: [][]byte{
			[]byte(`t2: {{ .t1 }}{{ .t1 }}
{{ .illegal }}: {{ printf "%s" .k2 }}
t3: {{ .t1 }}{{ .t1 }} additional text
normal: value`),
		},
		maxIterations: 1,
		res:           []byte{},
		err:           errors.New("illegal key found in line \"{{ .illegal }}: {{ printf \"%s\" .k2 }}\""),
	},
	{
		name: "illegal value",
		input: [][]byte{
			[]byte(`k1: v1
illegalValue: {`),
		},
		maxIterations: 1,
		res:           []byte{},
		err:           errors.New("yaml: line 2: did not find expected node content"),
	},
	{
		name: "incomplete rendering",
		input: [][]byte{
			[]byte(`t1: '{{ printf "%s" .k2 }}'
t2: '{{ .t1 }}{{ .t1 }}'
t3: '{{ .t1 }}{{ .t1 }} bla'
normal: value`),
		},
		maxIterations: 1,
		res:           []byte{},
		err:           errors.New("rendering incomplete after 1 iteration (max iterations 1)"),
	},
	{
		name: "merge & templating",
		input: [][]byte{
			[]byte(`k1: v1
k2: v2
k3:
  k31: v31
  k32:
    k321: v321
    k322: 
    - s3221
    - s3222
    k323: |
      bla
      blub`),
			[]byte(`k2: v2n
k3:
  k32:
    k321: 
    - k3211n
    - k3212n
    k323: |
      bla
      blub
      blee`),
			[]byte(`t1: '{{ printf "%s" .k2 }}'
t2: '{{ .t1 }}{{ .t1 }}'
t3: '{{ .t1 }}{{ .t1 }} bla'
normal: value`),
		},
		maxIterations: 5,
		res: []byte(`k1: v1
k2: v2n
k3:
  k31: v31
  k32:
    k321:
    - k3211n
    - k3212n
    k322:
    - s3221
    - s3222
    k323: |-
      bla
      blub
      blee
normal: value
t1: v2n
t2: v2nv2n
t3: v2nv2n bla
`),
		err: nil,
	},
}

func TestParseFiles(t *testing.T) {

	for _, test := range testValues {
		fmt.Printf("test %s\n", test.name)
		res, err := helmV.Render(test.input, test.maxIterations, true)
		if errDiff(test.err, err) {
			t.Error(errOutput(
				fmt.Sprintf("%s error", test.name),
				err,
				test.err,
			))
		}

		if len(res) != len(test.res) {
			fmt.Println(string(res))
			fmt.Println(string(test.res))
			t.Error(errOutput(
				fmt.Sprintf("%s result", test.name),
				string(res),
				string(test.res),
			))
			t.FailNow()
		}
		for i := range res {
			if res[i] != test.res[i] {
				t.Error(errOutput(
					fmt.Sprintf("%s result", test.name),
					string(res),
					string(test.res),
				))
			}
		}

	}
}

func errOutput(name string, result, expected interface{}) error {
	return fmt.Errorf("[MISMATCH] %s.\nResult: \"%s\" \nExpect: \"%s\"",
		name, result, expected)
}

func errDiff(errExp, errRes error) bool {
	if errExp == errRes {
		return false
	}
	if errExp == nil && errRes != nil {
		return true
	}
	if errExp != nil && errRes == nil {
		return true
	}
	if errExp.Error() != errRes.Error() {
		return true
	}
	return false
}
