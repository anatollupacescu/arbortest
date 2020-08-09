package runner

import (
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type node struct {
	ID     string `json:"id"`
	Status string `json:"status"`

	groupName string
}

type link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

type output struct {
	Commit  string `json:"commit"`
	Message string `json:"message"`
	Nodes   []node `json:"nodes"`
	Links   []link `json:"links"`

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

const (
	betweenTests  = 1
	betweenGroups = 3
)

func marshal(g Graph) string {
	statuses := []string{"skip", "fail", "pass"}

	out := output{
		Commit:  "unknown",
		Message: "unknown",
		nodes:   make(map[node]struct{}),
		links:   make(map[link]struct{}),
	}

	for i := range g.groups {
		grp := g.groups[i]
		groupNode := node{
			ID:     grp.name,
			Status: statuses[grp.status],
		}

		out.Node(groupNode)

		for _, tst := range grp.tests {
			testNode := node{
				ID:        tst.name,
				Status:    statuses[tst.status],
				groupName: grp.name,
			}
			out.Node(testNode)

			linkNode := link{
				Source: tst.name,
				Target: grp.name,
				Value:  betweenTests,
			}
			out.Link(linkNode)
		}
	}

	for fromGroupName := range g.deps {
		targetGroups := g.deps[fromGroupName]

		for _, destinationGroupName := range targetGroups {
			linkNode := link{
				Source: fromGroupName,
				Target: destinationGroupName,
				Value:  betweenGroups,
			}
			out.Link(linkNode)
		}
	}

	commit, message := g.infoProvider()
	out.Message = message
	out.Commit = commit

	data, err := json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}

func gitCommitAndMessage() (commit string, message string) {
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return
	}

	commit = string(out)
	commit = strings.TrimRight(commit, "\n")

	out, err = exec.Command("git", "show-branch", "--no-name", "HEAD").Output()
	if err != nil {
		log.Fatal("online ", err)
	}

	message = string(out)
	message = strings.TrimRight(message, "\n")

	return
}
