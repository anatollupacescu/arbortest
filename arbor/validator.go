package arbor

import (
	"fmt"
	"log"
)

// ErrInvalidProviderCall is returned when a tests calls a provider without an associated validator.
type ErrInvalidProviderCall struct {
	test, provider string
}

func (e *ErrInvalidProviderCall) Error() string {
	return fmt.Sprintf("%q calls invalid provider: %q", e.test, e.provider)
}

var (
	// ErrNoTestsDeclared returned when no testing code is found in the given test files.
	ErrNoTestsDeclared = fmt.Errorf("no tests found declared in the given folder files")
	// ErrIncompleteTestConfiguration returned when there is less than two providers declared
	ErrIncompleteTestConfiguration = fmt.Errorf("incomplete test configuration")
)

func validateErr(calls map[string][]string) (err error) {
	if len(calls) == 0 {
		return ErrNoTestsDeclared
	}

	var validProviders = validProviders(calls)

	for f, v := range calls {
		for i := range v {
			prov := v[i]
			if _, ok := validProviders[prov]; !ok {
				return &ErrInvalidProviderCall{
					f, v[i],
				}
			}
		}
	}

	return nil
}

func validProviders(calls map[string][]string) (vp map[string]bool) {
	vp = make(map[string]bool)

	for _, v := range calls {
		if len(v) == 1 {
			vp[v[0]] = true
		}
	}

	return
}

func printWarnings(callSuite suite) {
	var unusedProviders = unusedProviders(callSuite)

	for _, p := range unusedProviders {
		log.Printf("\u26a1 warning: %q declared but not used", p)
	}
}

func unusedProviders(callSuite suite) (unused []string) {
	var all = make([]string, 0, len(callSuite))

	for _, v := range callSuite {
		all = append(all, v...)
	}

	var providers []string

	for _, v := range callSuite {
		if len(v) == 1 {
			providers = append(providers, v...)
		}
	}

	for _, p := range providers {
		var found bool

		for _, e := range all {
			if p == e {
				found = true
				break
			}
		}

		if !found {
			unused = append(unused, p)
		}
	}

	return unused
}
