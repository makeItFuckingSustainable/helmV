package flatmap_test

import (
	"fmt"
	"hval/pkg/flatmap"
	"reflect"
	"testing"
)

var testFlatten = []struct {
	in  map[string]interface{}
	res map[string]flatmap.MapEntry
	err error
}{
	{
		in: map[string]interface{}{
			"fail": map[interface{}]interface{}{
				[2]string{"slice", "key"}: "val",
			},
		},
		res: map[string]flatmap.MapEntry{},
		err: fmt.Errorf("errors during flattening: (working partial key: \"%s\" - value: \"%s\"): cannot cast key \"%s\" to string,",
			[]string{"fail"}, "val", [2]string{"slice", "key"}),
	},
	{
		in: map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "val",
				}}},
		res: map[string]flatmap.MapEntry{
			"a.b.c": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "c"},
				Value:      "val",
			},
		},
		err: nil,
	},
	{
		in: map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "val",
					"d": []string{"res1", "res2"},
				}}},
		res: map[string]flatmap.MapEntry{
			"a.b.c": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "c"},
				Value:      "val",
			},
			"a.b.d": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "d"},
				Value:      []string{"res1", "res2"},
			}},
		err: nil,
	},
	{
		in: map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{
					"c": "val",
					"d": []string{"res1", "res2"},
				}},
			"x": map[string]interface{}{
				"b": map[string]interface{}{
					"d": map[string]interface{}{"hello": "map"},
				}},
		},
		res: map[string]flatmap.MapEntry{
			"a.b.c": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "c"},
				Value:      "val",
			},
			"a.b.d": flatmap.MapEntry{
				OrderedKey: []string{"a", "b", "d"},
				Value:      []string{"res1", "res2"},
			},
			"x.b.d.hello": flatmap.MapEntry{
				OrderedKey: []string{"x", "b", "d", "hello"},
				Value:      "map",
			},
		},
		err: nil,
	},
}

func TestFlatten(t *testing.T) {
	for _, test := range testFlatten {
		res, err := flatmap.Flatten(test.in)
		if err != test.err {
			if err.Error() != test.err.Error() {
				t.Errorf("expected error does not match.\nResult: %+v \nExpect: %+v", err, test.err)
			}
		}
		if !reflect.DeepEqual(res, test.res) {
			t.Errorf("maps are not equal.\nResult: %+v \nExpect: %+v", res, test.res)
		}

	}
}
