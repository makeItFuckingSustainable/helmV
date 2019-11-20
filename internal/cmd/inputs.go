package cmd

import (
	"bytes"
	"fmt"
	"github.com/makeItFuckingSustainable/helmV/internal/debug"
	"github.com/makeItFuckingSustainable/helmV/internal/render"
	"github.com/makeItFuckingSustainable/helmV/internal/yamltmpl"
	"github.com/makeItFuckingSustainable/helmV/pkg/flatmap"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func LoadInput(filePaths []string, d debug.Debugger) (
	flatmap.YamlMap,
	error,
) {
	aggMap := map[string]flatmap.MapEntry{}
	for _, path := range filePaths {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return flatmap.YamlMap{},
				fmt.Errorf("cannot read file \"%s\" - error: %s", f, err)
		}
		fSan, err := yamltmpl.Sanitize(f, d)
		if err != nil {
			return flatmap.YamlMap{}, err
		}
		y := map[string]interface{}{}
		if err := yaml.Unmarshal(fSan, &y); err != nil {
			return flatmap.YamlMap{}, err
		}

		fm, err := flatmap.Flatten(y)
		if err != nil {
			return flatmap.YamlMap{}, err
		}

		// merge results in flat map
		for k, v := range fm {
			aggMap[k] = v
		}
	}

	infl, err := flatmap.Inflate(aggMap)
	if err != nil {
		return flatmap.YamlMap{}, err
	}

	return infl, nil
}

func RenderResult(
	infl flatmap.YamlMap,
	output io.Writer,
	maxIterations uint,
	d debug.Debugger) error {
	inflBytes, err := yaml.Marshal(&infl)
	if err != nil {
		return err
	}
	tmpl, err := yamltmpl.Desanitize(inflBytes, d)
	if err != nil {
		return err
	}
	rendered := new(bytes.Buffer)
	err = render.Recursive(tmpl, infl, rendered, maxIterations)
	if err != nil {
		return err
	}
	var resParsed interface{}
	err = yaml.Unmarshal(rendered.Bytes(), &resParsed)
	if err != nil {
		return err
	}
	res, err := yaml.Marshal(resParsed)
	if err != nil {
		return err
	}
	_, err = output.Write(res)
	return err
}
