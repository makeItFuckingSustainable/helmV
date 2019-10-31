package flatmap_test

import (
	"fmt"
	"hval/pkg/flatmap"
	"reflect"
	"testing"
)

var testInflate = []struct {
	in       map[string]flatmap.MapEntry
	res      flatmap.YamlMap
	printRes string
	err      error
}{
	{
		in: map[string]flatmap.MapEntry{
			"fail": flatmap.MapEntry{
				OrderedKey: []string{},
				Value:      "val",
			}},
		res:      flatmap.YamlMap{},
		printRes: "{}",
		err:      fmt.Errorf("no key to insert provided for value %v", "val"),
	},
	{
		in: map[string]flatmap.MapEntry{
			"a.b.c": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "c"},
				Value:      "val",
			}},
		res: flatmap.YamlMap{
			"a": flatmap.YamlMap{
				"b": flatmap.YamlMap{
					"c": "val",
				}}},
		printRes: "{a: {b: {c: val}}}",
		err:      nil,
	},
	{
		in: map[string]flatmap.MapEntry{
			"a.b.c": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "c"},
				Value:      "val",
			},
			"a.b.d": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "d"},
				Value:      flatmap.YamlSlice{"res1", "res2"},
			}},
		res: flatmap.YamlMap{
			"a": flatmap.YamlMap{
				"b": flatmap.YamlMap{
					"c": "val",
					"d": flatmap.YamlSlice{"res1", "res2"},
				}}},
		printRes: "{a: {b: {d: [res1, res2], c: val}}}",
		err:      nil,
	},
	{
		in: map[string]flatmap.MapEntry{
			"a.b.c": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "c"},
				Value:      "val",
			},
			"a.b.d": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "d"},
				Value:      flatmap.YamlSlice{"res1", "res2"},
			},
			"x.b.d": flatmap.MapEntry{
				OrderedKey: []string{"x", "b", "d"},
				Value:      flatmap.YamlMap{"hello": "map"},
			},
		},
		res: flatmap.YamlMap{
			"a": flatmap.YamlMap{
				"b": flatmap.YamlMap{
					"c": "val",
					"d": flatmap.YamlSlice{"res1", "res2"},
				}},
			"x": flatmap.YamlMap{
				"b": flatmap.YamlMap{
					"d": flatmap.YamlMap{"hello": "map"},
				}},
		},
		printRes: "{a: {b: {c: val, d: [res1, res2]}}, x: {b: {d: {hello: map}}}}",
		err:      nil,
	},
}

func TestInflate(t *testing.T) {

	for _, test := range testInflate {
		res, err := flatmap.Inflate(test.in)
		if errDiff(test.err, err) {
			t.Errorf("expected error does not match.\nResult: %s \nExpect: %s", err, test.err)
		}

		if !reflect.DeepEqual(res, test.res) {
			t.Errorf("maps are not equal.\nResult: %+v \nExpect: %+v", res, test.res)
		}
	}
	t.FailNow()
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
