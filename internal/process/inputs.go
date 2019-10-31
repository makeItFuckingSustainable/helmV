package process

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func (i *input) LoadInput(filePaths []string) (
	map[string]map[string]interface{},
	error,
) {
	result := map[string]map[string]interface{}{}
	for _, path := range filePaths {
		fmt.Println(path)
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return map[string]map[string]interface{}{}, err
		}
		fSan, err := i.Sanitize(f)
		if err != nil {
			return map[string]map[string]interface{}{}, err
		}
		y := map[string]interface{}{}
		if err := yaml.Unmarshal(fSan, &y); err != nil {
			return map[string]map[string]interface{}{}, err
		}
		result[path] = y
	}
	return result, nil
}
