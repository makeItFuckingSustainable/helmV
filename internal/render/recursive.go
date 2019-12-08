package render

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

// Recursive accepts a template and input for that template and renders
// the template recursively with the input as values in each iteration.
// The recursion stops when either no change happens in the rendered template anymore
// or maxIterations is met.
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

	for t.hasChanged() && t.iteration < maxIterations {
		t.render()
	}

	if t.err != nil {
		return t.err
	}
	if t.hasChanged() {
		return fmt.Errorf(
			"rendering incomplete after %v iteration (max iterations %v)",
			t.iteration,
			maxIterations,
		)
	}
	_, err := output.Write(t.dataCurrent)
	return err
}

type tmpl struct {
	err          error
	dataCurrent  []byte
	dataPrevious []byte
	input        interface{}
	iteration    uint
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
