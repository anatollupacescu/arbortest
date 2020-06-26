package arbor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type (
	testName  string
	providers []string

	suite map[testName]providers
)

type ParseResult struct {
	Warnings []string
	Tests    suite
	Error    string
}

func Parse(src string) ParseResult {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)

	if err != nil {
		log.Fatalf("parse file: %v\n", err)
	}

	var calls = make(suite)
	var providerList providers

	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.FuncDecl); ok {
			test := gen.Name
			if strings.HasPrefix(test.Name, "provider") {
				providerList = append(providerList, test.Name)
				continue
			}
			pc := providerCalls(gen)
			calls[testName(test.Name)] = pc
		}
	}

	var warnings []string

	if len(calls) == 0 {
		warnings = append(warnings, "warning: no tests found")
	} else {
		var unusedProviders = getUnusedProviders(providerList, calls)

		for _, p := range unusedProviders {
			w := fmt.Sprintf("warning: '%s' declared but not used", p)
			warnings = append(warnings, w)
		}
	}

	var msg string

	if f, p := invalidProviders(providerList, calls); f != "" {
		msg = fmt.Sprintf("error: '%s' calls invalid provider: '%s'", f, p)
	}

	return ParseResult{
		Tests:    calls,
		Warnings: warnings,
		Error:    msg,
	}
}

func invalidProviders(pl providers, calls suite) (f testName, p string) {
	var valid = make(map[string]bool, 0)

	for _, v := range calls {
		if len(v) == 1 {
			valid[v[0]] = true
		}
	}

	for f, v := range calls {
		if len(v) > 1 {
			for i := range v {
				if _, ok := valid[v[i]]; !ok {
					return f, v[i]
				}
			}
		}
	}

	return
}

func getUnusedProviders(pl providers, calls suite) providers {
	var all providers

	for _, v := range calls {
		all = append(all, v...)
	}

	var unused providers

	for _, p := range pl {
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

func providerCalls(gen *ast.FuncDecl) (calls []string) {
	for _, s := range gen.Body.List {
		as, ok := s.(*ast.AssignStmt)

		if !ok {
			break
		}

		if len(as.Rhs) != 1 {
			break
		}

		rhs := as.Rhs[0]

		ce, ok := rhs.(*ast.CallExpr)

		if !ok {
			break
		}

		f, ok := ce.Fun.(*ast.Ident)

		if !ok {
			break
		}

		if strings.HasPrefix(f.Name, "provider") {
			calls = append(calls, f.Name)
		}
	}

	return
}
