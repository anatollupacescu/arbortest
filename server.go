package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	b := &broker{
		make(map[stringChan]struct{}),
		make(chan (stringChan)),
		make(chan (stringChan)),
		make(stringChan),
	}

	go b.listen()

	http.Handle("/events/", b)

	graphs := make([]string, 0)

	http.HandleFunc("/data/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if r.Method == "GET" {
			jsonData := strings.Join(graphs, ",")
			fmt.Fprintf(w, "[%s]", jsonData)

			return
		}

		if r.Method != "POST" {
			http.Error(w, "only GET and POST methods expected", http.StatusMethodNotAllowed)

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

		graph := string(bts)
		graphs = append([]string{graph}, graphs...)
		b.messages <- graph
	})

	// box := packr.New("demo", "./public")
	// dir := http.FileServer(box)
	// http.Handle("/", dir)

	portStr := fmt.Sprintf(":%d", *port)
	log.Printf("listening on port %s", portStr)

	return http.ListenAndServe(portStr, nil)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
