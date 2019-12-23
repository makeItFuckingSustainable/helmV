package helmV

import (
	"bufio"
	"bytes"
	"fmt"
	"log"

	"github.com/makeItFuckingSustainable/helmV/internal/render"
	"github.com/makeItFuckingSustainable/helmV/internal/yamltmpl"
	"github.com/makeItFuckingSustainable/helmV/pkg/flatmap"
	"gopkg.in/yaml.v2"
)

// Render is the single exported function of the helmV package. Given a slice of
// parsed input files, a maximal number of recursive iterations and a debug flag
// it renders a single resulting yaml file from all input files with all templates
// executed.
// If debug = true, the function will write debugging information to stdOut in
// of errors
func Render(files [][]byte, maxIt uint, debug bool) (
	[]byte, error,
) {
	infl, debugParse, err := parseFiles(files)
	if err != nil {
		if debug {
			writeDebug("parsing", debugParse)
		}
		return []byte{}, err
	}
	res, debugRender, err := renderTmpl(infl, maxIt)
	if err != nil {
		if debug {
			writeDebug("render", debugRender)
		}
		return []byte{}, err
	}
	return res, nil
}

// parseFiles takes the parsed input files, sanitizes their content and aggregates
// them into a single YamlMap which is returned.
func parseFiles(files [][]byte) (flatmap.YamlMap, []byte, error) {
	aggMap := map[string]flatmap.MapEntry{}
	for _, f := range files {
		fSan, err := yamltmpl.Sanitize(f)
		if err != nil {
			return flatmap.YamlMap{}, fSan, err
		}
		y := map[string]interface{}{}
		if err := yaml.Unmarshal(fSan, &y); err != nil {
			return flatmap.YamlMap{}, fSan, err
		}

		fm, err := flatmap.Flatten(y)
		if err != nil {
			return flatmap.YamlMap{}, []byte{}, err
		}

		// merge results in flat map
		for k, v := range fm {
			aggMap[k] = v
		}
	}

	infl, err := flatmap.Inflate(aggMap)
	if err != nil {
		return flatmap.YamlMap{}, []byte{}, err
	}

	return infl, []byte{}, nil
}

// Render is the main function of helmV. It takes a YamlMap as input and orchestrates
// the data preparation and execution of the recursive template rendering process.
func renderTmpl(infl flatmap.YamlMap, maxIterations uint) (
	[]byte, []byte, error,
) {
	inflBytes, err := yaml.Marshal(&infl)
	if err != nil {
		return []byte{}, nil, err
	}
	tmpl, err := yamltmpl.Desanitize(inflBytes)
	if err != nil {
		return []byte{}, tmpl, err
	}
	rendered := new(bytes.Buffer)
	err = render.Recursive(tmpl, infl, rendered, maxIterations)
	if err != nil {
		return []byte{}, rendered.Bytes(), err
	}
	var resParsed interface{}
	err = yaml.Unmarshal(rendered.Bytes(), &resParsed)
	if err != nil {
		return []byte{}, rendered.Bytes(), err
	}
	res, err := yaml.Marshal(resParsed)
	return res, []byte{}, err
}

func writeDebug(label string, out []byte) {
	fmt.Println("")
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		log.Printf("DEBUG (%s) | %s\n", label, scanner.Text())
	}
}
