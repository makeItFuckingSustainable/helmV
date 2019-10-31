package flatmap

import (
	"fmt"
	"strings"
)

// YamlMap implements a map that returns a map in valid yaml format when printed
type YamlMap map[string]interface{}

func (m YamlMap) String() string {
	res := make([]string, 0, len(m))
	for k, v := range m {
		res = append(res, fmt.Sprintf("%s: %s", k, v))
	}
	return fmt.Sprintf("{%s}", strings.Join(res, ", "))
}

// YamlSlice implements a slice that returns a slice in valid yaml format when printed
type YamlSlice []interface{}

func (s YamlSlice) String() string {
	res := make([]string, 0, len(s))
	for _, v := range s {
		res = append(res, fmt.Sprintf("%s", v))
	}
	return fmt.Sprintf("[%s]", strings.Join(res, ", "))
}

// Inflate takes a flattened map (result of topmap.Flatten function) and returns
// the map to the original nested form (inflates it)
func Inflate(m map[string]MapEntry) (result YamlMap, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("flatmap: %v", r)
				return
			}
			err = fmt.Errorf("%v [recovered]", err)
		}
	}()

	result = make(YamlMap)
	for _, v := range m {
		if len(v.OrderedKey) == 0 {
			result = YamlMap{}
			err = fmt.Errorf("no key to insert provided for value %v", v.Value)
			return
		}
		result = upsert(result, v.OrderedKey, v.Value)
	}
	return
}

func upsert(
	m interface{}, orderedKeys []string, value interface{},
) YamlMap {
	switch cast := m.(type) {
	case YamlMap:
		if subMap, ok := cast[orderedKeys[0]]; ok {
			subres := upsert(subMap, orderedKeys[1:], value)
			cast[orderedKeys[0]] = subres
			return cast
		}
		cast[orderedKeys[0]] = overwrite(orderedKeys[1:], value)
		return cast
	default:
		panic(fmt.Errorf("you should not get here - https://xkcd.com/2200/"))
	}
}

func overwrite(orderedKeys []string, value interface{}) interface{} {
	if len(orderedKeys) == 0 {
		return castValue(value)
	}
	// build the remaining submap from leaf element orderedKeys[max]
	// to root element orderedKeys[0]
	value = castValue(value)
	for i := len(orderedKeys) - 1; i > 0; i-- {
		value = YamlMap{orderedKeys[i]: value}
	}
	return YamlMap{orderedKeys[0]: value}
}

func castValue(value interface{}) interface{} {
	switch cast := value.(type) {
	case []interface{}:
		return YamlSlice(cast)
	default:
		return cast
	}
}
