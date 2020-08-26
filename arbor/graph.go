package arbor

import (
	"container/list"
	"errors"
	"fmt"
)

type testDescriptor struct {
	Name, Title string
}

type testGroup struct {
	Deps  []string
	Tests []testDescriptor
}

type graph struct {
	order  []string
	groups map[string]testGroup
}

var errRepeatedDeclaration = errors.New("repeated declaration of 'after' for current group")

func (g graph) addGroup(id string, deps []string, testName testDescriptor) error {
	group := g.groups[id]

	if len(group.Deps) > 0 && len(deps) > 0 {
		return errRepeatedDeclaration
	}

	group.Deps = append(group.Deps, deps...)
	group.Tests = append(group.Tests, testName)

	g.groups[id] = group

	return nil
}

func (g *graph) saveOrder() {
	initialOrder := list.New()

	for _, v := range g.order {
		initialOrder.PushBack(v)
	}

	ordered := list.New()

	for e := initialOrder.Front(); e != nil; e = e.Next() {
		group := e.Value

		if !dependenciesAreSatisfied(g, group.(string), ordered) {
			initialOrder.PushBack(group)
			continue
		}

		ordered.PushBack(group)
	}

	g.order = make([]string, 0, initialOrder.Len())
	for e := ordered.Front(); e != nil; e = e.Next() {
		g.order = append(g.order, e.Value.(string))
	}
}

func dependenciesAreSatisfied(g *graph, group string, left *list.List) bool {
	for _, dep := range g.groups[group].Deps {
		var found bool

		for e := left.Front(); e != nil; e = e.Next() {
			if e.Value.(string) == dep {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

var errCircularDependency = errors.New("circular dependency")

func (g graph) circularLinks() error {
	for _, id := range g.order {
		group := g.groups[id]
		for _, dep := range group.Deps {
			if hasLinkTo(g, dep, id) {
				return fmt.Errorf("%w %s->%s->%s", errCircularDependency, id, dep, id)
			}
		}
	}

	return nil
}

var errGroupNotFound = errors.New("group not found")

func (g graph) missingGroups() error {
	for _, groupID := range g.order {
		grp := g.groups[groupID]
		for _, depID := range grp.Deps {
			if _, hasGroup := g.groups[depID]; !hasGroup {
				return fmt.Errorf("%w: %s", errGroupNotFound, depID)
			}
		}
	}

	return nil
}

func (g graph) isEmpty() bool {
	return len(g.groups) == 0
}

func hasLinkTo(g graph, from, to string) bool {
	target := g.groups[from]
	for _, dep := range target.Deps {
		if dep == to {
			return true
		}
	}

	return false
}
