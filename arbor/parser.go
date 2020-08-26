package arbor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type testBundle struct {
	comment, testTitle, testName string
}

func parse(src string) []testBundle {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, "src.go", src, parser.ParseComments)

	if err != nil {
		log.Fatalf("parse file: %v\n", err)
	}

	var bundles []testBundle

	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.FuncDecl); ok && hasTestSignature(gen) {
			comment := getComment(gen)
			if comment == "" {
				continue
			}

			testName := gen.Name.Name
			testTitle := strings.TrimPrefix(testName, "test")

			bundles = append(bundles, testBundle{
				comment:   comment,
				testTitle: testTitle,
				testName:  testName,
			})
		}
	}

	return bundles
}

func getComment(gen *ast.FuncDecl) string {
	comments := gen.Doc.List
	if len(comments) != 1 {
		return ""
	}

	comment := comments[0]

	return comment.Text
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

	param := params.List[0]

	return hasCorrectParameterName(param)
}

func hasCorrectParameterName(field *ast.Field) bool {
	if ident, ok := field.Type.(*ast.StarExpr); ok {
		if selector, ok := ident.X.(*ast.SelectorExpr); ok {
			if selectorX, ok := selector.X.(*ast.Ident); ok && selectorX.Name != "runner" {
				return false
			}

			return selector.Sel.Name == "T"
		}
	}

	return false
}
