package main

import (
	"bytes"
	"fmt"
	"hval/internal/process"
	"hval/internal/render"
	"hval/pkg/flags"
	"hval/pkg/flatmap"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func aggregate(files map[string]map[string]interface{}) (
	map[string]flatmap.MapEntry,
	error,
) {
	aggregatedValues := map[string]flatmap.MapEntry{}
	for _, content := range files {
		fm, err := flatmap.Flatten(content)
		if err != nil {
			return map[string]flatmap.MapEntry{}, err
		}
		for k, v := range fm {
			aggregatedValues[k] = v
		}
	}
	return aggregatedValues, nil
}

func main() {

	args, err := flags.Parse()
	check(err)
	debug := new(bytes.Buffer)
	yamlProcess := process.New(debug, false)
	files, err := yamlProcess.LoadInput(args.Files)
	check(err)
	aggMap, err := aggregate(files)
	check(err)
	infl, err := flatmap.Inflate(aggMap)
	check(err)
	inflBytes, err := yaml.Marshal(&infl)
	tmpl, err := yamlProcess.Desanitize(inflBytes)
	output := new(bytes.Buffer)
	check(render.Recursive(tmpl, infl, output, 10))
	var out interface{}
	check(yaml.Unmarshal(output.Bytes(), &out))
	res, err := yaml.Marshal(out)
	check(err)
	fmt.Println(string(res))
	check(ioutil.WriteFile("output.yaml", res, 0644))

}

func check(err error) {
	// TODO add proper error and log handling
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}
