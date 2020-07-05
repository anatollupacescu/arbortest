package runner

import "encoding/json"

type node struct {
	ID     string `json:"id"`
	Group  int    `json:"group"`
	Status string `json:"status"`
}

type link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

type output struct {
	Nodes []node `json:"nodes"`
	Links []link `json:"links"`
}

const defaultWeight = 3

func marshal(tests ...*test) string {
	out := output{
		Nodes: make([]node, 0),
		Links: make([]link, 0),
	}

	for _, t := range tests {
		marshalTest(t, &out)
	}

	bytes, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func marshalTest(ts *test, out *output) {
	for _, t := range ts.deps {
		marshalTest(t, out)
	}

	for _, d := range ts.deps {
		l := link{
			Source: ts.name,
			Target: d.name,
			Value:  defaultWeight, //TODO make configurable
		}

		out.Links = append(out.Links, l)
	}

	if contains(out.Nodes, ts.name) {
		return
	}

	var statuses = []string{"pending", "fail", "pass"}

	n := node{
		ID:     ts.name,
		Group:  int(ts.status),
		Status: statuses[ts.status],
	}

	out.Nodes = append(out.Nodes, n)
}

func contains(ns []node, name string) bool {
	for i := 0; i < len(ns); i++ {
		if ns[i].ID == name {
			return true
		}
	}

	return false
}
