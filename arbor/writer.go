package arbor

import (
	"fmt"
	"sort"
	"strings"
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
		%s,
	}`

var testsSrc = `tests := map[string]func() error{
		%s,
	}`

var dependenciesSrc = `dependencies := map[string][]string{
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
	var providers []string

	for i := range d.providers {
		str := fmt.Sprintf("\"%s\"", d.providers[i])
		providers = append(providers, str)
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i] < providers[j]
	})

	commaSep := strings.Join(providers, ", ")

	return fmt.Sprintf("\"%s\": {%s}", d.testName, commaSep)
}

func GenerateSource(pkg string, s suite) string {
	var (
		validatorList  []string
		dependencyList []string
		testList       []string
	)

	for testName, providers := range s {
		if len(providers) == 1 {
			vr := validator{
				testName: string(testName),
				provider: providers[0],
			}
			validatorList = append(validatorList, vr.String())
		} else {
			dep := dependency{
				testName:  string(testName),
				providers: []string(providers),
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

	vals := fmt.Sprintf(validatorSrc, strings.Join(validatorList, ", "))
	tests := fmt.Sprintf(testsSrc, strings.Join(testList, ", "))
	deps := fmt.Sprintf(dependenciesSrc, strings.Join(dependencyList, ", "))

	return fmt.Sprintf(main, pkg, vals, deps, tests)
}
