package loader

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type loader struct {
	n int
	c int

	requests []*http.Request
	ws       *websocket.Conn
}

func New(n, c int, ws *websocket.Conn) *loader {
	requests := make([]*http.Request, 0)
	return &loader{n, c, requests, ws}
}

func (l *loader) AddTasks(method string, host string, path string) {
	// let's just start with one
	// TODO: this is going to be a loop
	req, err := http.NewRequest(method, host+"/"+path, nil)

	if err != nil {
		fmt.Println("something went wrong !")
	}

	l.requests = append(l.requests, req)
}

func (l *loader) spawnClient(wg *sync.WaitGroup, queue chan []*http.Request) {
	// this is one out of c clients
	client := &http.Client{}

	// if we're not busy, get a tasks list from the queue channel
	for requests := range queue {
		for _, req := range requests {
			// Send the request
			res, err := client.Do(req)
			if err != nil {
				fmt.Println("wut?")
			}

			res.Body.Close()

		}

		if err := l.ws.WriteMessage(websocket.TextMessage, []byte("One tasks batch done")); err != nil {
			fmt.Println("An error occured: the requested test doesn't exist")
		}
		// We are done
		wg.Done()
	}
}

func (l *loader) Run() {
	// We'll block until all requests are done
	var wg sync.WaitGroup
	wg.Add(l.n)

	// The queue channel will contain an array of tasks (up to n time)
	queue := make(chan []*http.Request, l.n)

	// Spawn all our clients
	for i := 0; i < l.c; i++ {
		go func() {
			l.spawnClient(&wg, queue)
		}()
	}

	// Populate our channel with tasks arrays
	for i := 0; i < l.n; i++ {
		queue <- l.requests
	}

	// Close our channel
	close(queue)

	// unblock
	wg.Wait()
}
