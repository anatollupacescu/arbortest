package runner

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//nolint:gochecknoglobals	//idiomatic way of working with flags in Go
var uri = flag.String("arborURL", "", "arbor server url")

// Upload sends the graph json to the UI server.
func Upload(data string) {
	flag.Parse()

	if *uri == "" {
		log.Println("server uri not provided, skipping upload...")
		return
	}

	u, err := url.ParseRequestURI(*uri)
	if err != nil {
		log.Fatalf("validate upload uri: %s", err)
	}

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(data))

	if err != nil {
		log.Fatalf("build upload request: %s", err)
	}

	defer func() {
		_ = req.Body.Close()
	}()

	hc := http.Client{}

	r, err := hc.Do(req)
	if err != nil {
		log.Fatalf("upload test results: %s", err)
	}

	defer func() {
		_ = r.Body.Close()
	}()

	if r.StatusCode != http.StatusOK {
		log.Fatalf("upload test results: bad response status: %s", r.Status)
		return
	}

	log.Println("upload test results: success")
}
