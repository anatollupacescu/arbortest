package arbor

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

func build(bundles []testBundle) (out graph, err error) {
	grph := graph{
		groups: make(map[string]testGroup),
	}

	err = populate(grph, bundles)
	if err != nil {
		return
	}

	for k := range grph.groups {
		grph.order = append(grph.order, k)
	}

	// sort alphabetically for consistent errors
	sort.Slice(grph.order, func(i, j int) bool {
		return grph.order[i] < grph.order[j]
	})

	if err := grph.missingGroups(); err != nil {
		return out, err
	}

	if err := grph.circularLinks(); err != nil {
		return out, err
	}

	grph.saveOrder()

	return grph, nil
}

func populate(g graph, bundles []testBundle) (err error) {
	for i := range bundles {
		input := bundles[i]
		if err = applyBundle(g, input); err != nil {
			err = fmt.Errorf("%w near '%s':\n%s", err, input.testName, input.comment)
			return
		}
	}

	return nil
}

var (
	errMissingGroupDeclaration = errors.New("missing group declaration")
	errDuplicateToken          = errors.New("duplicate token")
	errUnexpectedSegmentKind   = errors.New("unexpected kind")
)

func applyBundle(g graph, bundle testBundle) error {
	c := bundle.comment
	c = strings.TrimLeft(c, "/")
	c = strings.Trim(c, " ")

	var (
		groupID       string
		afterDeclared bool
		dependencies  []string
	)

	inputs := strings.Split(c, " ")
	for _, input := range inputs {
		seg, err := newFromString(input)
		if err != nil {
			return err
		}

		switch seg.kind {
		case "group":
			if groupID != "" {
				return errDuplicateToken
			}

			groupID = seg.groupID
		case "after":
			if afterDeclared {
				return errDuplicateToken
			}

			afterDeclared = true
			dependencies = seg.dependencies
		default:
			return errUnexpectedSegmentKind
		}
	}

	if groupID == "" {
		return errMissingGroupDeclaration
	}

	testDesc := testDescriptor{
		Name:  bundle.testName,
		Title: bundle.testTitle,
	}

	if err := g.addGroup(groupID, dependencies, testDesc); err != nil {
		return err
	}

	return nil
}

type segment struct {
	kind         string
	groupID      string
	dependencies []string
}

const keyValuePairSize = 2

var errBadToken = errors.New("bad token")

func newFromString(in string) (segment, error) {
	var seg segment

	components := strings.Split(in, ":")

	if len(components) != keyValuePairSize {
		return seg, errBadToken
	}

	kind := components[0]
	switch kind {
	case "group":
		groupID, err := extractGroupID(components[1])
		if err != nil {
			return seg, err
		}

		seg.kind = kind
		seg.groupID = groupID
	case "after":
		dependencies, err := extractDependencies(components[1])
		if err != nil {
			return seg, err
		}

		seg.kind = kind
		seg.dependencies = dependencies
	default:
		return seg, errBadToken
	}

	return seg, nil
}

func extractGroupID(in string) (string, error) {
	if strings.Contains(in, ",") {
		return "", errBadToken
	}

	return in, nil
}

var errEmptyValueNotAllowed = errors.New("empty value not allowed")

func extractDependencies(component string) (elements []string, err error) {
	if component == "" {
		return nil, errEmptyValueNotAllowed
	}

	elements = strings.Split(component, ",")

	for i := range elements {
		elem := elements[i]
		if elem == "" {
			return nil, errEmptyValueNotAllowed
		}
	}

	return elements, nil
}
