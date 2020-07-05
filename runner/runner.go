package runner

import "testing"

type status int

const (
	pending = status(iota)
	fail
	pass
)

type test struct {
	name    string
	deps    []*test
	runFunc func(t *testing.T)
	status  status
}

type (
	dependencies map[string][]string
	tests        map[string]func(t *testing.T)
)

// Run entry point.
func Run(t *testing.T, d dependencies, tt tests) string {
	newTests := make([]*test, 0, len(tt))

	for name, testFunc := range tt {
		providers := d[name]
		if len(providers) == 1 {
			continue
		}

		newTests = append(newTests, &test{
			name:    name,
			runFunc: testFunc,
			deps:    computeDeps(providers, d, tt),
		})
	}

	for _, v := range newTests {
		v := v
		t.Run(v.name, func(t2 *testing.T) {
			v.run(t2)
		})
	}

	return marshal(newTests...)
}

func computeDeps(providers []string, d dependencies, t tests) []*test {
	out := make([]*test, 0)

	for _, givenProvider := range providers {
		for testName, testProviders := range d {
			if len(testProviders) != 1 {
				continue
			}

			if testProviders[0] == givenProvider {
				out = append(out, &test{
					name:    testName,
					runFunc: t[testName],
				})
			}
		}
	}

	return out
}

func (ts *test) run(t *testing.T) {
	if hasFailingDependencies(ts.deps, t) {
		return
	}

	ts.runFunc(t)
}

func hasFailingDependencies(deps []*test, t *testing.T) bool {
	if len(deps) == 0 {
		return false
	}

	for _, dep := range deps {
		switch dep.status {
		case pass:
			continue
		case fail:
			return true
		case pending:
			dep.run(t)

			if t.Failed() {
				return true
			}
		default:
			return false
		}
	}

	return false
}
