package arbor

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

type testResult struct {
	Output string
	Error  string
}

func Run(validators map[string]string, dependencies map[string][]string, tests map[string]func() error) testResult {
	var validProviders []string

	var (
		links []link
		nodes []test
	)

	for name, fun := range tests {
		providerName, isProvider := validators[name]

		if !isProvider {
			continue
		}

		links = append(links, link{
			source: name,
			target: providerName,
		})

		node := test{
			name:   name,
			status: passed,
		}

		if err := fun(); err != nil {
			node.status = failed
			node.errMsg = err.Error()
		}

		if node.status == passed {
			validProviders = append(validProviders, providerName)
		}

		nodes = append(nodes, node)
		nodes = append(nodes, test{
			name:   providerName,
			status: node.status,
		})
	}

	//execution
	for name, fun := range tests {
		_, isProvider := validators[name]

		if isProvider {
			continue
		}

		dependsOn := dependencies[name]

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

	return testResult{
		Output: marshall(nodes, links),
		Error:  nodesError(nodes),
	}
}

func nodesError(tests []test) string {
	for _, t := range tests {
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

func Upload(s interface{}) string {
	return "uploaded"
}
