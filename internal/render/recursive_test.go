package render_test

import (
	"bytes"
	"fmt"
	"helmV/internal/render"
	"helmV/pkg/flatmap"
	"testing"
)

var testValues = []struct {
	name   string
	tmpl   string
	values interface{}
	maxIt  uint
	res    string
	err    error
}{
	{
		name: "no templating",
		tmpl: `
k2: world
normal: value
`,
		values: flatmap.YamlMap{},
		maxIt:  10,
		res: `
k2: world
normal: value
`,
		err: nil,
	},
	{
		name: "recursive templating",
		tmpl: `
k2: world
t1: {{ printf "hello-%s" .k2 }}
t2: {{ .t1 }}{{ .t1 }}
t3: {{ .t2 }} bla
`,
		values: flatmap.YamlMap{
			"k2": "world",
			"t1": `{{ printf "hello-%s" .k2 }}`,
			"t2": `{{ .t1 }}{{ .t1 }}`,
			"t3": `{{ .t2 }} bla`,
		},
		maxIt: 10,
		res: `
k2: world
t1: hello-world
t2: hello-worldhello-world
t3: hello-worldhello-world bla
`,
		err: nil,
	},
	{
		name: "nested map insertion",
		tmpl: `
v4:
	v41: test
t4: {{ .v4 }}
`,
		values: flatmap.YamlMap{
			"v4": flatmap.YamlMap{"v41": "test"},
			"t4": `{{ .v4 }}`,
		},
		maxIt: 10,
		res: `
v4:
	v41: test
t4: {v41: test}
`,
		err: nil,
	},
	{
		name: "recursive nested map insertion",
		tmpl: `
v4:
	v41: test
t4: {{ .v4 }}
t5:
	t51: {{ .t4 }}
t6: {{ .t5.t51 }}
normal: value
`,
		values: flatmap.YamlMap{
			"v4":     flatmap.YamlMap{"v41": "test"},
			"t4":     `{{ .v4 }}`,
			"t5":     flatmap.YamlMap{"t51": `{{ .t4 }}`},
			"normal": "value",
		},
		maxIt: 10,
		res: `
v4:
	v41: test
t4: {v41: test}
t5:
	t51: {v41: test}
t6: {v41: test}
normal: value
`,
		err: nil,
	},
	{
		name: "nested slice insertion",
		tmpl: `
v5:
	v51:
	- s1
	- s2
t7: {{ .v5 }}
`,
		values: flatmap.YamlMap{
			"v5": flatmap.YamlMap{"v51": flatmap.YamlSlice{"s1", "s2"}},
		},
		maxIt: 10,
		res: `
v5:
	v51:
	- s1
	- s2
t7: {v51: [s1, s2]}
`,
		err: nil,
	},
	{
		name: "partial templating",
		tmpl: `
k2: world
t1: {{ printf "hello-%s" .k2 }}
t2: {{ .t1 }}{{ .t1 }}
t3: {{ .t2 }} bla
`,
		values: flatmap.YamlMap{
			"k2": "world",
			"t1": `{{ printf "hello-%s" .k2 }}`,
			"t2": `{{ .t1 }}{{ .t1 }}`,
			"t3": `{{ .t2 }} bla`,
		},
		maxIt: 1,
		res:   ``,
		err:   fmt.Errorf("rendering incomplete after 1 iteration (max iterations 1)"),
	},
	{
		name: "invalid values",
		tmpl: `
k2: world
t1: {{ printf "hello-%s" .k2 }}
`,
		values: "invalid input",
		maxIt:  1,
		res:    ``,
		err:    fmt.Errorf("template: tmp_0:3:25: executing \"tmp_0\" at <.k2>: can't evaluate field k2 in type string"),
	},
	{
		name: "invalid template",
		tmpl: `
k2: world
t1: {{ printf "hello-%s" .k2
`,
		values: flatmap.YamlMap{},
		maxIt:  1,
		res:    ``,
		err:    fmt.Errorf("template: tmp_0:3: unclosed action"),
	},
}

func TestRender(t *testing.T) {

	for _, test := range testValues {
		fmt.Printf("test %s\n", test.name)
		res := new(bytes.Buffer)
		err := render.Recursive(
			[]byte(test.tmpl),
			test.values,
			res,
			test.maxIt,
		)
		if errDiff(test.err, err) {
			t.Error(errOutput(
				fmt.Sprintf("%s error", test.name),
				err,
				test.err,
			))
		}

		if res.String() != test.res {
			t.Error(errOutput(
				fmt.Sprintf("%s result", test.name),
				res.String(),
				test.res,
			))
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
