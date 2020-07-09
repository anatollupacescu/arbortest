package runner

import (
	"sort"
	"testing"
)

const (
	pending = iota
	fail
	pass
	skip
)

type test struct {
	name      string
	providers []string
	runFunc   func(t *testing.T)
	status    int
}

type (
	dependencies map[string][]string
	tests        map[string]func(t *testing.T)
)

// Run entry point.
func Run(t *testing.T, dd dependencies, tt tests) string {
	all := make([]*test, 0, len(tt))

	for testName, providers := range dd {
		all = append(all, &test{
			name:      testName,
			providers: providers,
			runFunc:   tt[testName],
			status:    pass,
		})
	}

	sort.Slice(all, func(i, j int) bool {
		return len(all[i].providers) < len(all[j].providers)
	})

	failedProviders := make(map[string]bool)

	for _, a := range all {
		var hasAFailedProvider bool

		for _, ap := range a.providers {
			if _, found := failedProviders[ap]; found {
				hasAFailedProvider = true
				break
			}
		}

		if hasAFailedProvider {
			a.status = skip
			continue
		}

		a := a
		t.Run(a.name, func(t *testing.T) {
			a.runFunc(t)
			if t.Failed() {
				a.status = fail
			}

			if a.status == fail && len(a.providers) == 1 {
				failedProviders[a.providers[0]] = true
			}
		})
	}

	return marshal(all...)
}
