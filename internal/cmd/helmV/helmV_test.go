package helmV_test

import (
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
		// 		res: flatmap.YamlMap{
		// 			"k1": "v1",
		// 			"k2": "v2n",
		// 			"k3": flatmap.YamlMap{
		// 				"k31": "v31",
		// 				"k32": flatmap.YamlMap{
		// 					"k321": flatmap.YamlSlice{"k3211n", "k3212n"},
		// 					"k322": flatmap.YamlSlice{"s3221", "s3222"},
		// 					"k323": `bla
		// blub
		// blee`,
		// 				}},
		// 			"normal": "value",
		// 			"t1":     `{{ printf "%s" .k2 }}`,
		// 			"t2":     `{{ .t1 }}{{ .t1 }}`,
		// 			"t3":     `{{ .t1 }}{{ .t1 }} bla`,
		// 		},
	},
}

func TestParseFiles(t *testing.T) {

	for _, test := range testValues {
		fmt.Printf("test %s\n", test.name)
		res, err := helmV.Render(test.input, test.maxIterations, false)
		if errDiff(test.err, err) {
			t.Error(errOutput(
				fmt.Sprintf("%s error", test.name),
				err,
				test.err,
			))
		}

		if len(res) != len(test.res) {
			fmt.Println(res)
			fmt.Println(test.res)
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
