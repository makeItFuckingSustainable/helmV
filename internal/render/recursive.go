package render

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

type tmpl struct {
	err          error
	dataCurrent  []byte
	dataPrevious []byte
	input        interface{}
	iteration    uint
}

func NewTemplate(
	template []byte,
	input interface{},
	output io.Writer,
) (tmpl, error) {

	// // TODO: implement optional check for cyclic references
	// // potentially check on (flatInput map[string]flatmap.MapEntry)
	// if cyclicReferences(flatInput) {
	// 	return tmpl{}, fmt.Errorf("cyclic dependency in aggregated values")
	// }

	return tmpl{
		err:          nil,
		dataCurrent:  template,
		dataPrevious: []byte{},
		input:        input,
		iteration:    0,
	}, nil
}

func Recursive(
	template []byte,
	input interface{},
	output io.Writer,
	maxIterations uint,
) error {
	t := tmpl{
		err:          nil,
		dataCurrent:  template,
		dataPrevious: []byte{},
		input:        input,
		iteration:    0,
	}

	for t.hasChanged() || t.iteration < maxIterations {
		t.render()
		if t.err != nil {
			return t.err
		}
	}
	_, err := output.Write(t.dataCurrent)
	return err
}

func (t *tmpl) hasChanged() bool {
	if n := bytes.Compare(t.dataPrevious, t.dataCurrent); n == 0 {
		return false
	}
	return true
}

func (t *tmpl) render() {
	if t.err != nil {
		return
	}
	t.dataPrevious = make([]byte, len(t.dataCurrent))
	copy(t.dataPrevious, t.dataCurrent)

	tt, err := template.New(
		fmt.Sprintf("tmp_%d", t.iteration),
	).Parse(string(t.dataPrevious))
	if err != nil {
		t.err = err
		return
	}
	res := new(bytes.Buffer)

	if err := tt.Execute(res, t.input); err != nil {
		t.err = err
		return
	}
	t.iteration++
	t.dataCurrent = res.Bytes()
	return
}
