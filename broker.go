package main

import (
	"fmt"
	"net/http"
)

type stringChan chan string

type broker struct {
	clients        map[stringChan]struct{}
	newClients     chan stringChan
	defunctClients chan stringChan
	messages       stringChan
}

func (b *broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			b.clients[s] = struct{}{}
		case s := <-b.defunctClients:
			delete(b.clients, s)
			close(s)
		case msg := <-b.messages:
			for s := range b.clients {
				s <- msg
			}
		}
	}
}

func (b *broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	messageChan := make(stringChan)

	b.newClients <- messageChan

	go func() {
		<-r.Context().Done()
		b.defunctClients <- messageChan
	}()

	enableCors(&w)
	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	for {
		msg, open := <-messageChan

		if !open {
			break
		}

		fmt.Fprintf(w, "data: %s\n\n", msg)
		f.Flush()
	}
}
