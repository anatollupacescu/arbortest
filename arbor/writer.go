package arbor

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
)

const tmpl = `package {{ .Package }}

import (
	"testing"

	arbor "github.com/anatollupacescu/arbortest/runner"
)

func TestArbor(t *testing.T) {
	g := arbor.New()
{{ $groups := .Groups }}{{ range $elem := .Order }}
	t.Run({{printf "%q" $elem}}, func(t *testing.T) {
		at := arbor.NewT(t)
		g.Group({{ printf "%q" $elem }}){{ $testGroup := (index $groups $elem)}}{{ $len := (len $testGroup.Deps) }}{{ if (gt $len 0) }}
		g.After(at, {{ $testGroup.Deps | commaSep }}){{end}}{{ range $test := $testGroup.Tests}}
		g.Append(at, {{ printf "%q" $test }}, {{ $test }}){{ end }}
	})
{{ end }}
	output := g.JSON()

	arbor.Upload(output)
}
`

func commaSep(elems []string) string {
	quoted := make([]string, 0, len(elems))
	for _, e := range elems {
		quoted = append(quoted, fmt.Sprintf("%q", e))
	}

	return strings.Join(quoted, ",")
}

func generateSource(pkg string, g graph) string {
	data := struct {
		Package string
		Order   []string
		Groups  map[string]testGroup
	}{
		Package: pkg,
		Order:   g.order,
		Groups:  g.groups,
	}

	fmap := template.FuncMap{
		"commaSep": commaSep,
	}

	parse, err := template.New("test").Funcs(fmap).Parse(tmpl)

	tmpl := template.Must(parse, err)

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, data)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}
