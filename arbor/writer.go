package arbor

import (
	"fmt"
	"sort"
	"strings"
)

const main = `package %s

import (
	"testing"

	"github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	%s
	%s
	%s
	if r := runner.Run(validators, dependencies, tests); r.Error != "" {
		t.Error(r.Error)
	}
}
`

const validatorSrc = `validators := map[string]string{
		%s,
	}`

const testsSrc = `tests := map[string]func() error{
		%s,
	}`

const dependenciesSrc = `dependencies := map[string][]string{
		%s,
	}`

type validator struct {
	testName, provider string
}

func (v validator) String() string {
	return fmt.Sprintf("\"%s\": \"%s\"", v.testName, v.provider)
}

type dependency struct {
	testName  string
	providers []string
}

func (d dependency) String() string {
	var providers []string = make([]string, len(d.providers))

	for i := range d.providers {
		str := fmt.Sprintf("\"%s\"", d.providers[i])
		providers[i] = str
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i] < providers[j]
	})

	commaSep := strings.Join(providers, ", ")

	return fmt.Sprintf("\"%s\": {%s}", d.testName, commaSep)
}

func generateSource(pkg string, calls map[string][]string) string {
	var (
		validatorList  []string
		dependencyList []string
	)

	testList := make([]string, 0, len(calls))

	for testName, providers := range calls {
		if len(providers) == 1 {
			vr := validator{
				testName: testName,
				provider: providers[0],
			}
			validatorList = append(validatorList, vr.String())
		} else {
			dep := dependency{
				testName:  testName,
				providers: providers,
			}
			dependencyList = append(dependencyList, dep.String())
		}

		str := fmt.Sprintf("\"%s\": %s", testName, testName)
		testList = append(testList, str)
	}

	sort.Slice(validatorList, func(i, j int) bool {
		return validatorList[i] < validatorList[j]
	})
	sort.Slice(testList, func(i, j int) bool {
		return testList[i] < testList[j]
	})
	sort.Slice(dependencyList, func(i, j int) bool {
		return dependencyList[i] < dependencyList[j]
	})

	var vals, tests, deps string

	if len(validatorList) > 0 {
		vals = fmt.Sprintf(validatorSrc, strings.Join(validatorList, ", "))
	} else {
		vals = "var validators map[string]string"
	}

	if len(dependencyList) > 0 {
		deps = fmt.Sprintf(dependenciesSrc, strings.Join(dependencyList, ", "))
	} else {
		deps = "var dependencies map[string][]string"
		vals = "var validators map[string]string"
	}

	tests = fmt.Sprintf(testsSrc, strings.Join(testList, ", "))

	return fmt.Sprintf(main, pkg, vals, deps, tests)
}
