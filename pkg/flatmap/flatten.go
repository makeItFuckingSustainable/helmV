package flatmap

import (
	"fmt"
	"strings"
)

type resError struct {
	Result MapEntry
	Error  error
}

// MapEntry holds the full information about a nested key-value pair. The OrderedKey
// field holds the ordered keys which identify the value in a nested map
// (outer -> inner as left -> right)
type MapEntry struct {
	OrderedKey []string
	Value      interface{}
}

// Flatten takes in a nested map and transforms it into a flat map where each key
// is a concatenation of the nested keys. The value of the flat map holds the value
// and an ordered list of the original nested keys. E.g.
// map["k1":map["k2":"v"]] -> map["k1.k2":MapEntry{[]string{"k1", "k2"}, "v"}]
// This function returns an error if it finds keys that cannot be casted to strings.
func Flatten(m map[string]interface{}) (map[string]MapEntry, error) {
	openBranches := make(chan int, 1)
	flatKV := make(chan resError, 1)
	openBranches <- +1
	go sendKV(m, []string{}, flatKV, openBranches)
	resMap, errors, err := aggregateKV(flatKV, openBranches)

	if err != nil {
		errMsg := "errors during flattening:"
		if len(errors) > 0 {
			for _, kv := range errors {
				errMsg = fmt.Sprintf("%s (working partial key: \"%s\" - value: \"%s\"): %s,",
					errMsg, kv.Result.OrderedKey, kv.Result.Value, kv.Error)
			}
		}
		return map[string]MapEntry{}, fmt.Errorf(errMsg)
	}
	return resMap, nil
}

func sendKV(
	m interface{}, orderedKeys []string,
	flatKV chan<- resError, openBranches chan<- int,
) {
	switch cast := m.(type) {
	case map[string]interface{}:
		for k, v := range cast {
			openBranches <- +1
			go sendKV(v, append(orderedKeys, k), flatKV, openBranches)
		}
	case map[interface{}]interface{}:
		for k, v := range cast {
			ks, ok := k.(string)
			if !ok {
				flatKV <- resError{
					Result: MapEntry{OrderedKey: orderedKeys, Value: v},
					Error:  fmt.Errorf("cannot cast key \"%s\" to string", k),
				}
				break
			}
			openBranches <- +1
			go sendKV(v, append(orderedKeys, ks), flatKV, openBranches)
		}
	default:
		openBranches <- +1
		flatKV <- resError{
			Result: MapEntry{OrderedKey: orderedKeys, Value: cast},
			Error:  nil,
		}
	}
	openBranches <- -1
}

func aggregateKV(flatKV <-chan resError, openBranches chan int) (
	map[string]MapEntry, []resError, error,
) {
	doneCh := make(chan bool)
	var done bool
	go areWeDone(openBranches, doneCh)

	resMap := make(map[string]MapEntry)
	errors := []resError{}
	for !done {
		select {
		case done = <-doneCh:
		case kve := <-flatKV:
			openBranches <- -1
			if kve.Error != nil {
				errors = append(errors, kve)
				break
			}
			resMap[strings.Join(kve.Result.OrderedKey, ".")] = kve.Result
		}
	}
	if len(errors) != 0 {
		return map[string]MapEntry{}, errors, fmt.Errorf("failed to flatten input")
	}
	return resMap, []resError{}, nil
}

func areWeDone(openBranches <-chan int, done chan<- bool) {
	nbBranches := 0
	nbBranches += <-openBranches
	for nbBranches > 0 {
		nbBranches += <-openBranches
	}
	done <- true
	close(done)
}
