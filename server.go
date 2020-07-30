package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gobuffalo/packr/v2"
)

func init() {
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("» ")
}

//nolint:gochecknoglobals // idiomatic way of working with flags in Go
var port = flag.Int("port", 3000, "port to listen to")

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	http.HandleFunc("/data.json", dataFunc)

	box := packr.New("demo", "./public")
	dir := http.FileServer(box)
	http.Handle("/", dir)

	portStr := fmt.Sprintf(":%d", *port)
	log.Printf("listening on port %s", portStr)

	return http.ListenAndServe(portStr, nil)
}

//nolint:gochecknoglobals	//simplest way to persist graph data
var data = `{
	"nodes": [{"id": "No tests ran yet", "status": "pending"}],
	"links": []
}`

func dataFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, data)

		return
	}

	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	data = string(bts)
}
