package helmV

import (
	"bytes"

	"github.com/makeItFuckingSustainable/helmV/internal/render"
	"github.com/makeItFuckingSustainable/helmV/internal/yamltmpl"
	"github.com/makeItFuckingSustainable/helmV/pkg/flatmap"
	"github.com/makeItFuckingSustainable/helmV/pkg/logerrs"
	"gopkg.in/yaml.v2"
)

func ParseFiles(files [][]byte, d logerrs.Debugger) (
	flatmap.YamlMap,
	error,
) {
	aggMap := map[string]flatmap.MapEntry{}
	for _, f := range files {
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

func Render(
	infl flatmap.YamlMap,
	maxIterations uint,
	d logerrs.Debugger) ([]byte, error) {
	inflBytes, err := yaml.Marshal(&infl)
	if err != nil {
		return []byte{}, err
	}
	tmpl, err := yamltmpl.Desanitize(inflBytes, d)
	if err != nil {
		return []byte{}, err
	}
	rendered := new(bytes.Buffer)
	err = render.Recursive(tmpl, infl, rendered, maxIterations)
	if err != nil {
		return []byte{}, err
	}
	var resParsed interface{}
	err = yaml.Unmarshal(rendered.Bytes(), &resParsed)
	if err != nil {
		return []byte{}, err
	}
	return yaml.Marshal(resParsed)
}
