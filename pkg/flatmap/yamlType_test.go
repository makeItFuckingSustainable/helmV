package flatmap_test

import (
	"helmV/pkg/flatmap"
	"testing"
)

var testYamlTypes = []struct {
	in  interface{}
	res string
	err error
}{
	{
		in:  flatmap.YamlMap{},
		res: "{}",
	},
	{
		in:  flatmap.YamlMap{"a": "b"},
		res: "{a: b}",
	},
	{
		in: flatmap.YamlMap{
			"a": flatmap.YamlMap{
				"b": flatmap.YamlMap{
					"c": "val",
				}}},
		res: "{a: {b: {c: val}}}",
	},
	{
		in:  flatmap.YamlSlice{},
		res: "[]",
	},
	{
		in:  flatmap.YamlSlice{"a", "b"},
		res: "[a, b]",
	},
	{
		in: flatmap.YamlSlice{
			"a",
			flatmap.YamlMap{
				"b": flatmap.YamlMap{
					"c": "val",
				}},
			"c",
		},
		res: "[a, {b: {c: val}}, c]",
	},
}

func TestYamlTypes(t *testing.T) {

	for _, test := range testYamlTypes {
		switch cast := test.in.(type) {
		case flatmap.YamlMap:
			if cast.String() != test.res {
				t.Errorf(
					"printed %T does not match.\nResult: %+v \nExpect: %+v",
					cast, cast.String(), test.res,
				)
			}
		case flatmap.YamlSlice:
			if cast.String() != test.res {
				t.Errorf(
					"printed %T does not match.\nResult: %+v \nExpect: %+v",
					cast, cast.String(), test.res,
				)
			}
		default:
			t.Errorf("unexpected test input %T: %+v", cast, cast)
		}
	}
}
