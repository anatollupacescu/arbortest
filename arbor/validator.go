package arbor

import (
	"fmt"
)

// ErrInvalidProviderCall is returned when a tests calls a provider without an associated validator.
type ErrInvalidProviderCall struct {
	test, provider string
}

func (e *ErrInvalidProviderCall) Error() string {
	return fmt.Sprintf("%q calls invalid provider: %q", e.test, e.provider)
}

// ErrNoTestsDeclared returned when no testing code is found in the given test files.
var ErrNoTestsDeclared = fmt.Errorf("no tests found declared in the given folder files")

func validateErr(calls map[string][]string) (errors []error) {
	if len(calls) == 0 {
		errors = append(errors, ErrNoTestsDeclared)
		return
	}

	var validProviders = validProviders(calls)

	for f, v := range calls {
		for i := range v {
			prov := v[i]
			if _, ok := validProviders[prov]; !ok {
				e := ErrInvalidProviderCall{
					f, v[i],
				}
				errors = append(errors, &e)
			}
		}
	}

	return
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

func validateWarns(callSuite suite) (warnings []string) {
	if len(callSuite.calls) == 0 {
		warnings = append(warnings, "no tests found")
		return
	}

	var unusedProviders = unusedProviders(callSuite)

	for _, p := range unusedProviders {
		w := fmt.Sprintf("%q declared but not used", p)
		warnings = append(warnings, w)
	}

	return
}

func unusedProviders(callSuite suite) (unused []string) {
	var all = make([]string, 0, len(callSuite.calls))

	for _, v := range callSuite.calls {
		all = append(all, v...)
	}

	for _, p := range callSuite.providers {
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
