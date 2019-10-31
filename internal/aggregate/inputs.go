package aggregate

import (
	"fmt"
	"hval/pkg/flatmap"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Inputs(files []string) (map[string]flatmap.MapEntry, error) {
	result := map[string]flatmap.MapEntry{}
	for _, file := range files {
		y, err := readYaml(file)
		if err != nil {
			return map[string]flatmap.MapEntry{}, err
		}
		fm, err := flatmap.Flatten(y)
		if err != nil {
			return map[string]flatmap.MapEntry{}, err
		}
		for k, v := range fm {
			result[k] = v
		}
	}
	return result, nil
}

func readYaml(f string) (map[string]interface{}, error) {
	yamlFile, err := ioutil.ReadFile(f)
	if err != nil {
		return map[string]interface{}{},
			fmt.Errorf("cannot read file \"%s\" - error: %s", f, err)
	}
	result := map[string]interface{}{}
	err = yaml.Unmarshal(yamlFile, result)
	if err != nil {
		return map[string]interface{}{},
			fmt.Errorf("cannot parse file \"%s\" - error: %s", f, err)
	}
	return result, nil
}
