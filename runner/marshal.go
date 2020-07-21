package runner

import (
	"encoding/json"
	"log"
)

type node struct {
	ID     string `json:"id"`
	Group  status `json:"group"`
	Status string `json:"status"`

	groupName string
}

type link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

type output struct {
	Nodes []node `json:"nodes"`
	Links []link `json:"links"`

	nodes map[node]struct{}
	links map[link]struct{}
}

func (o *output) Node(n node) {
	if _, isPresent := o.nodes[n]; isPresent {
		return
	}

	o.Nodes = append(o.Nodes, n)

	o.nodes[n] = struct{}{}
}

func (o *output) Link(l link) {
	if _, isPresent := o.links[l]; isPresent {
		return
	}

	o.Links = append(o.Links, l)

	o.links[l] = struct{}{}
}

const defaultLinkNodeValue = 3

func marshal(g Graph) string {
	statuses := []string{"skip", "fail", "pass"}

	out := output{
		nodes: make(map[node]struct{}),
		links: make(map[link]struct{}),
	}

	for i := range g.groups {
		grp := g.groups[i]
		groupNode := node{
			ID:     grp.name,
			Group:  grp.status,
			Status: statuses[grp.status],
		}

		out.Node(groupNode)

		for _, tst := range grp.tests {
			testNode := node{
				ID:        tst.name,
				Group:     tst.status,
				groupName: grp.name,
				Status:    statuses[tst.status],
			}
			out.Node(testNode)

			linkNode := link{
				Source: tst.name,
				Target: grp.name,
				Value:  defaultLinkNodeValue,
			}
			out.Link(linkNode)
		}
	}

	for fromGroupName := range g.deps {
		targetGroups := g.deps[fromGroupName]

		source := fromGroupName + "-ext"
		groupStatus := g.groups.get(fromGroupName).status

		testNode := node{
			ID:     source,
			Group:  groupStatus,
			Status: statuses[groupStatus],
		}
		out.Node(testNode)

		linkNode := link{
			Source: source,
			Target: fromGroupName,
			Value:  defaultLinkNodeValue,
		}
		out.Link(linkNode)

		for _, destinationGroupName := range targetGroups {
			destination := destinationGroupName + "-ext"
			groupStatus := g.groups.get(fromGroupName).status
			testNode := node{
				ID:     destination,
				Group:  groupStatus,
				Status: statuses[groupStatus],
			}
			out.Node(testNode)

			linkNode := link{
				Source: source,
				Target: destination,
				Value:  defaultLinkNodeValue,
			}
			out.Link(linkNode)

			extToHome := link{
				Source: destination,
				Target: destinationGroupName,
				Value:  defaultLinkNodeValue,
			}
			out.Link(extToHome)
		}
	}

	data, err := json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}
