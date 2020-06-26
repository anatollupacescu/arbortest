package arbor

import (
	"encoding/json"
	"log"
	"sort"
)

type outNode struct {
	ID     string `json:"id"`
	Group  int    `json:"group"`
	Status string `json:"status"`
}

type outLink struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

func marshall(tests []test, lnks []link) string {
	var (
		nodes []outNode
		links []outLink
	)

	if tests == nil {
		nodes = make([]outNode, 0)
	}

	if lnks == nil {
		links = make([]outLink, 0)
	}

	for _, t := range tests {
		var status = "skipped"
		switch t.status {
		case failed:
			status = "failed"
		case passed:
			status = "passed"
		}
		nodes = append(nodes, outNode{
			ID:     t.name,
			Group:  int(t.status),
			Status: status,
		})
	}

	for _, l := range lnks {
		links = append(links, outLink{
			Source: l.source,
			Target: l.target,
			Value:  0,
		})
	}

	var output = struct {
		Nodes []outNode `json:"nodes"`
		Links []outLink `json:"links"`
	}{
		Nodes: sortByID(nodes),
		Links: sortBySource(links),
	}

	str, err := json.Marshal(output)

	if err != nil {
		log.Fatalf("marshal output: %s", err)
	}

	return string(str)
}

func sortByID(nodes []outNode) []outNode {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].ID < nodes[j].ID
	})

	return nodes
}

func sortBySource(links []outLink) []outLink {
	sort.Slice(links, func(i, j int) bool {
		return links[i].Source < links[j].Source
	})

	return links
}
