package runner

type status uint8

const (
	skip status = iota
	fail
	pass
)

type group struct {
	name   string
	status status
	tests  []test
}

type test struct {
	name   string
	status status
}

type groups []group

func (gg *groups) add(g group) {
	*gg = append(*gg, g)
}

func (gg groups) get(name string) *group {
	for i := range gg {
		g := &gg[i]
		if g.name == name {
			return g
		}
	}

	return new(group)
}

// Graph exported.
type Graph struct {
	groups           groups
	deps             map[string][]string
	currentGroupName string
}

// New exported.
func New() *Graph {
	return &Graph{
		groups: make(groups, 0),
		deps:   make(map[string][]string),
	}
}

// After exported.
func (g Graph) After(t *T, dependencies ...string) {
	g.deps[g.currentGroupName] = dependencies

	for _, dependsOn := range dependencies {
		dep := g.groups.get(dependsOn)
		if dep.status != pass {
			t.Errorf("skipping '%s' because dependency '%s' has failed", g.currentGroupName, dependsOn)
			g.groups.get(g.currentGroupName).status = skip

			return
		}
	}
}

// Group exported.
func (g *Graph) Group(name string) {
	g.currentGroupName = name
	g.groups.add(group{
		name:   name,
		status: pass,
		tests:  make([]test, 0),
	})
}

// Append exported.
func (g *Graph) Append(t *T, name string, f func(t *T)) {
	grp := g.groups.get(g.currentGroupName)
	if t.Failed() || grp.status == skip {
		grp.tests = append(grp.tests, test{
			name:   name,
			status: skip,
		})

		return
	}

	f(t)

	node := test{
		name:   name,
		status: pass,
	}

	if t.Failed() {
		grp.status = fail
		node.status = fail
	}

	grp.tests = append(grp.tests, node)
}

// JSON exported.
func (g Graph) JSON() string {
	return marshal(g)
}
