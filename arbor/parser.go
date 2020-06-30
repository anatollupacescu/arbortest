package arbor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type suite struct {
	calls     map[string][]string
	providers []string
}

func parseSource(src string) suite {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)

	if err != nil {
		log.Fatalf("parse file: %v\n", err)
	}

	var (
		calls     = make(map[string][]string)
		providers []string
	)

	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.FuncDecl); ok {
			test := gen.Name

			if strings.HasPrefix(test.Name, "provider") {
				providers = append(providers, test.Name)
				continue
			}

			if strings.HasPrefix(test.Name, "test") {
				pc := providerCalls(gen)
				calls[test.Name] = pc
			}
		}
	}

	return suite{
		calls:     calls,
		providers: providers,
	}
}

func providerCalls(gen *ast.FuncDecl) (calls []string) {
	for _, s := range gen.Body.List {
		f, ok := fromDecl(s)
		if ok && strings.HasPrefix(f.Name, "provider") {
			calls = append(calls, f.Name)
		}
	}

	return
}

func fromDecl(s interface{}) (f *ast.Ident, ok bool) {
	as, ok := s.(*ast.AssignStmt)

	if !ok || len(as.Rhs) != 1 {
		return nil, false
	}

	ce, ok := as.Rhs[0].(*ast.CallExpr)

	if !ok {
		return nil, false
	}

	if f, ok = ce.Fun.(*ast.Ident); !ok {
		return nil, false
	}

	return f, true
}
