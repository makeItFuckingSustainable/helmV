[![Coverage Status](https://coveralls.io/repos/github/makeItFuckingSustainable/helmV/badge.svg?branch=master)](https://coveralls.io/github/makeItFuckingSustainable/helmV?branch=master)
[![Build Status](https://travis-ci.com/makeItFuckingSustainable/helmV.svg?branch=master)](https://travis-ci.com/makeItFuckingSustainable/helmV)
[![Go Report Card](https://goreportcard.com/badge/github.com/makeItFuckingSustainable/helmV)](https://goreportcard.com/report/github.com/makeItFuckingSustainable/helmV)

# helmV

This project makes a naive version of go templating available within values of yaml
files. In addition it allows to deep-merge multiple yaml files into one single file.

The general process is that first all inputs are merged into a single `flatMap`.
This `flatMap` is subsequently used as input to render all remaining go templates
that occur in the `flatMap`. The process stops when either no change in the `flatMap`
occurs anymore, or the maximal number of iterations is exceeded.

## License

The yaml package is licensed under the Apache License 2.0. Please see the LICENSE file for details.

## Limitations

Naive version of go templating means, that only yaml templates can be processed
whose yaml keys are not templates themselves, see the following yaml template example:

```yaml
allowedKey: allowedValue
alsoAllowed: {{ .allowedKey }}{{ .allowedKey }}
{{ .illegal }}: {{ printf "%s" .k2 }}
```

## Build CLI

helmV can be used as a standalone cli. To build the cli binary run

```bash
export helmV_folder=<choose source code folder (e.g. ${HOME}/src/github.com/helmV)>
export branch=master
export name=helmV

git clone git@github.com:makeItFuckingSustainable/helmV.git \
    --branch ${branch} \
    --single-branch \
    --depth=1 \
    ${helmV_folder}
go test "${helmV_folder}/..." -v -race
go build -o ${name} ${helmV_folder}
chmod +x "${helmV_folder}/${name}"
```

## Usage

### As package

You can use helmV directly in your go project as a package. Here is an example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/makeItFuckingSustainable/helmV/internal/cmd/helmV"
)

func main() {

	merge1 := []byte(`
k1: v1
k2: v2
k3:
  k31: v31
  k32:
    k321: v321
    k322:
    - s3221
    - s3222
    k323: |
      bla
      blub
`)

	merge2 := []byte(`
k2: v2n
k3:
  k32:
    k321:
    - k3211n
    - k3212n
    k323: |
      bla
      blee
    k333: v333
`)

	template := []byte(`
t1: '{{ printf "%s" .k3.k32.k333 }}'
t2: '{{ .t1 }}{{ .t1 }}'
t3: '{{ .t1 }}{{ .t1 }} bla'
`)

	maxIterations := uint(1)
	debugging := false

	merged, err := helmV.Render(
        [][]byte{merge1, merge2},
        maxIterations,
        debugging,
    )
	if err != nil {
		log.Fatalf("[ERROR] helmV: %v", err)
	}
	fmt.Printf("--- merge dump:\n%s\n\n", string(merged))

	maxIterations = 5
	mergeTempl, err := helmV.Render(
        [][]byte{merge1, merge2, template},
        maxIterations,
        debugging,
    )
	if err != nil {
		log.Fatalf("[ERROR] helmV: %v", err)
	}
	fmt.Printf("--- merge-templated dump:\n%s\n\n", string(mergeTempl))
}

```

### CLI

To use the cli, first you have to [build](#build-cli) the helmV binary and either
make it available in your `${PATH}` or call it directly from its location.

The interface for helmV has the form

```shell
Usage of ./helmV:
  -debug
        Activate debugging.
  -max-iterations uint
        Maximal number of recursive iterations that helmV will execute. (default 10)
  -output string
        Absolute output path. Will default to "${PWD}/values.yaml". (default "values.yaml")
  -v value
        File holding input values. Relative path will be changed to absolute path as "${PWD}/filename". Multiple value files are processed first to last.
  -values value
        File holding input values. Relative path will be changed to absolute path as "${PWD}/filename". Multiple value files are processed first to last.
```

`N` yaml files can then be merged and contained templates rendered as

```go
helmV -max-iterations 5 -output res.yaml \
    -v file_1.yaml \
    -v file_2.yaml \
    ...
    -v file_N.yaml \
```

Here, the parameter `-max-iterations 5`specifies that the recursive template
rendering happens at most 5 times.
