package arbor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	tt := []struct {
		name     string
		inputs   []testBundle
		expected graph
		err      string
	}{
		{
			name: "unexpected keyword",
			inputs: []testBundle{
				{testName: "test1", comment: "// zebra:a"},
			},
			err: "bad token in\n{// zebra:a test1}",
		}, {
			name: "empty value for 'after' declaration",
			inputs: []testBundle{
				{testName: "test1", comment: "// after:"},
			},
			err: "empty value not allowed in\n{// after: test1}",
		}, {
			name: "bad value for 'after' declaration",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:one after:,"},
			},
			err: "empty value not allowed in\n{// group:one after:, test1}",
		}, {
			name: "single group with a misplaced character",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a,"},
			},
			err: "bad token in\n{// group:a, test1}",
		}, {
			name: "repeated 'group' declaration in the same comment",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a group:b"},
			},
			err: "duplicate token in\n{// group:a group:b test1}",
		}, {
			name: "repeated 'after' declaration in the same comment",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a after:b after:c"},
			},
			err: "duplicate token in\n{// group:a after:b after:c test1}",
		}, {
			name: "missing 'group' declaration in comment",
			inputs: []testBundle{
				{testName: "test1", comment: "// after:b"},
			},
			err: "missing group declaration in\n{// after:b test1}",
		}, {
			name: "duplicate group definition bundles tests together",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a"},
				{testName: "test2", comment: "// group:a"},
			},
			expected: graph{
				order: []string{"a"},
				groups: map[string]testGroup{
					"a": {
						Tests: []string{"test1", "test2"},
					},
				},
			},
		}, {
			name: "single group",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a"},
			},
			expected: graph{
				order: []string{"a"},
				groups: map[string]testGroup{
					"a": {
						Tests: []string{"test1"},
					},
				},
			},
		}, {
			name: "single group with dependency on non existent group",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a after:b"},
			},
			err: "group not found: b",
		}, {
			name: "two groups, one depending on another",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a"},
				{testName: "test2", comment: "// group:b after:a"},
			},
			expected: graph{
				order: []string{"a", "b"},
				groups: map[string]testGroup{
					"a": {
						Tests: []string{"test1"},
					},
					"b": {
						Tests: []string{"test2"},
						Deps:  []string{"a"},
					},
				},
			},
		}, {
			name: "two groups, one depending on another, order not important",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a after:b"},
				{testName: "test2", comment: "// group:b"},
			},
			expected: graph{
				order: []string{"b", "a"},
				groups: map[string]testGroup{
					"a": {
						Tests: []string{"test1"},
						Deps:  []string{"b"},
					},
					"b": {
						Tests: []string{"test2"},
					},
				},
			},
		}, {
			name: "three groups with dependencies",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a"},
				{testName: "test2", comment: "// group:b"},
				{testName: "test3", comment: "// group:c after:a,b"},
			},
			expected: graph{
				order: []string{"a", "b", "c"},
				groups: map[string]testGroup{
					"a": {
						Tests: []string{"test1"},
					},
					"b": {
						Tests: []string{"test2"},
					},
					"c": {
						Tests: []string{"test3"},
						Deps:  []string{"a", "b"},
					},
				},
			},
		}, {
			name: "three groups with dependencies on the first one",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a"},
				{testName: "test2", comment: "// group:b after:a"},
				{testName: "test3", comment: "// group:c after:a,b"},
			},
			expected: graph{
				order: []string{"a", "b", "c"},
				groups: map[string]testGroup{
					"a": {
						Tests: []string{"test1"},
					},
					"b": {
						Tests: []string{"test2"},
						Deps:  []string{"a"},
					},
					"c": {
						Tests: []string{"test3"},
						Deps:  []string{"a", "b"},
					},
				},
			},
		}, {
			name: "rejects circular dependency",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a after:b"},
				{testName: "test2", comment: "// group:b after:a"},
			},
			err: "circular dependency a->b->a",
		}, {
			name: "rejects repeated 'after' declaration within the same group",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a after:b"},
				{testName: "test2", comment: "// group:a after:c"},
				{testName: "test3", comment: "// group:b"},
				{testName: "test4", comment: "// group:c"},
			},
			err: "repeated declaration of 'after' for current group in\n{// group:a after:c test2}",
		}, {
			name: "orders by dependencies",
			inputs: []testBundle{
				{testName: "test1", comment: "// group:a after:b,c"},
				{testName: "test2", comment: "// group:a"},
				{testName: "test3", comment: "// group:b after:c"},
				{testName: "test4", comment: "// group:c"},
			},
			expected: graph{
				order: []string{"c", "b", "a"},
				groups: map[string]testGroup{
					"a": {
						Tests: []string{"test1", "test2"},
						Deps:  []string{"b", "c"},
					},
					"b": {
						Tests: []string{"test3"},
						Deps:  []string{"c"},
					},
					"c": {
						Tests: []string{"test4"},
					},
				},
			},
		},
	}

	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			real, err := build(tst.inputs)
			assert.Equal(t, tst.expected, real)
			if tst.err != "" {
				assert.EqualError(t, err, tst.err)
			}
		})
	}
}
