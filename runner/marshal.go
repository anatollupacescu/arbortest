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
	var statuses = []string{"pending", "fail", "pass", "skip"}

	out := output{
		Nodes: make([]node, 0),
		Links: make([]link, 0),
	}

	providers := make(map[string]node)

	for _, ts := range tests {
		n := node{
			ID:     ts.name,
			Group:  ts.status,
			Status: statuses[ts.status],
		}

		out.Nodes = append(out.Nodes, n)

		for _, providerName := range ts.providers {
			providers[providerName] = node{
				ID:     providerName,
				Group:  ts.status,
				Status: statuses[ts.status],
			}

			l := link{
				Source: ts.name,
				Target: providerName,
				Value:  defaultWeight, //TODO make configurable
			}

			out.Links = append(out.Links, l)
		}
	}

	for i := range providers {
		out.Nodes = append(out.Nodes, providers[i])
	}

	bytes, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
