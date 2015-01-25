package loader

import (
	"fmt"
	"net/http"
	"sync"
)

type loader struct {
	n int
	c int

	requests []*http.Request
}

func New(n, c int) *loader {
	requests := make([]*http.Request, 0)
	return &loader{n, c, requests}
}

func (l *loader) AddTasks() {
	// let's just start with one
	// TODO: this is going to be a loop
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080", nil)

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
		// Just do the first request for now (TODO: loop over all of them)
		req := requests[0]

		// Send the request
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("wut?")
		}

		res.Body.Close()

		// We are done
		wg.Done()

		fmt.Println("i finished a request")
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
