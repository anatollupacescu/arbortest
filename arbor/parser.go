package arbor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type suite map[string][]string

func parseSource(src string) suite {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)

	if err != nil {
		log.Fatalf("parse file: %v\n", err)
	}

	calls := make(map[string][]string)

	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.FuncDecl); ok && hasTestSignature(gen) {
			test := gen.Name.Name
			calls[test] = providerCalls(gen)
		}
	}

	return suite(calls)
}

func hasTestSignature(gen *ast.FuncDecl) bool {
	if !strings.HasPrefix(gen.Name.Name, "test") {
		return false
	}

	if gen.Type.Results.NumFields() > 0 {
		return false
	}

	params := gen.Type.Params
	if params.NumFields() != 1 {
		return false
	}

	p := params.List[0]
	if ident, ok := p.Type.(*ast.StarExpr); ok {
		if selector, ok := ident.X.(*ast.SelectorExpr); ok {
			if selectorX, ok := selector.X.(*ast.Ident); ok && selectorX.Name != "testing" {
				return false
			}

			return selector.Sel.Name == "T"
		}
	}

	return false
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
