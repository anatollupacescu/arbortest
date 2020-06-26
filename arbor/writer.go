package arbor

import (
	"fmt"
	"io/ioutil"
	"log"
)

var main = `package %s

import (
	"testing"

	"github.com/anatollupacescu/arbortest/arbor"
)

func TestArbor(t *testing.T) {
	%s
	%s
	%s
	if r := arbor.Run(validators, dependencies, tests); r.Error != "" {
		t.Error(r.Error)
	}
}
`

var validatorSrc = `validators := map[string]string{
		%s
	}`

var testsSrc = `tests := map[string]func() error{
		%s
	}`

var dependenciesSrc = `dependencies := map[string][]string{
		%s
	}`

func GenerateSource(pkg string, s suite) string {
	var (
		validatorList  string
		testList       string
		dependencyList string
	)

	for testName, providers := range s {
		if len(providers) == 1 {
			validatorList += fmt.Sprintf("\"%s\": \"%s\", ", testName, providers[0])
		} else {
			dependencyList += toDependencyList(testName, providers)
		}
		testList += fmt.Sprintf("\"%s\": %s, ", testName, testName)
	}

	vals := fmt.Sprintf(validatorSrc, validatorList)
	tests := fmt.Sprintf(testsSrc, testList)
	deps := fmt.Sprintf(dependenciesSrc, dependencyList)

	return fmt.Sprintf(main, pkg, vals, deps, tests)
}

func toDependencyList(testName testName, vals []string) (out string) {
	var commaSep string
	for i := range vals {
		commaSep += fmt.Sprintf("\"%s\", ", vals[i])
	}
	return fmt.Sprintf("\"%s\": {%s}, ", testName, commaSep)
}

func Write(fileName string, contents string) {
	err := ioutil.WriteFile(fileName, []byte(contents), 0644)

	if err != nil {
		log.Fatalf("create generated test file: %v", err)
	}
}
