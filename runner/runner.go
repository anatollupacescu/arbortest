package runner

type status int

const (
	skipped status = iota
	passed
	failed
)

type link struct {
	source, target string
}

type test struct {
	name   string
	status status
	errMsg string
}

// Result holds the serialized output of the successful run or the reason of failure.
type Result struct {
	Output string
	Error  string
}

type (
	validators   map[string]string
	dependencies map[string][]string
	tests        map[string]func() error
)

// Run takes a collection of validators, dependencies and tests and runs them in a prioritized order.
func Run(v validators, d dependencies, t tests) Result {
	nodes, links, validProviders := validateProviders(v, t)

	//execution
	for name, fun := range t {
		if _, isProvider := v[name]; isProvider {
			continue
		}

		dependsOn := d[name]

		for _, dep := range dependsOn {
			links = append(links, link{
				source: name,
				target: dep,
			})
		}

		node := test{
			name:   name,
			status: passed,
		}

		canBeRun := allDepsAreValid(dependsOn, validProviders)

		if !canBeRun {
			node.status = skipped
			node.errMsg = "failed dependency"
		} else if err := fun(); err != nil {
			node.status = failed
			node.errMsg = err.Error()
		}

		nodes = append(nodes, node)
	}

	return Result{
		Output: marshall(nodes, links),
		Error:  nodesError(nodes),
	}
}

func validateProviders(v validators, t tests) ([]test, []link, []string) {
	nodes := make([]test, 0, len(t))
	links := make([]link, 0, len(t))

	uniqProviders := make(map[string][]string)
	for v, p := range v {
		uniqProviders[p] = append(uniqProviders[p], v)
	}

	orderedByName := make([]string, 0, len(uniqProviders))
	for providerFunctionName := range uniqProviders {
		orderedByName = append(orderedByName, providerFunctionName)
	}

	validProviders := make([]string, 0, len(uniqProviders))

	for _, providerFunctionName := range orderedByName {
		validatorList := uniqProviders[providerFunctionName]

		var hasFailed bool

		for _, validationTestName := range validatorList {
			links = append(links, link{
				source: validationTestName,
				target: providerFunctionName,
			})

			node := test{
				name:   validationTestName,
				status: skipped,
			}

			if !hasFailed {
				fun := t[validationTestName]
				if err := fun(); err != nil {
					node.status = failed
					node.errMsg = err.Error()
					hasFailed = true
				} else {
					node.status = passed
				}
			}

			nodes = append(nodes, node)
		}

		providerStatus := failed

		if !hasFailed {
			validProviders = append(validProviders, providerFunctionName)
			providerStatus = passed
		}

		nodes = append(nodes, test{
			name:   providerFunctionName,
			status: providerStatus,
		})
	}

	return nodes, links, validProviders
}

func nodesError(tests []test) string {
	for i := range tests {
		t := tests[i]
		if t.status == failed {
			return t.errMsg
		}
	}

	return ""
}

func allDepsAreValid(deps, all []string) bool {
	if len(all) < len(deps) {
		return false
	}

	for _, dep := range deps {
		var found bool

		for _, vDep := range all {
			if dep == vDep {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// func upload(s interface{}) string {
// 	return "uploaded"
// }
