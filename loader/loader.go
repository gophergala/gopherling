package loader

import (
	"bytes"
	"fmt"
	"github.com/gophergala/gopherling/models"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type feedback struct {
	Task       int           `json:"task"`
	StatusCode int           `json:"statusCode"`
	Duration   time.Duration `json:"duration"`
}

type Task struct {
	method  string
	host    string
	path    string
	headers []models.Header
	rawBody string
}

func (t *Task) Request() *http.Request {
	body := bytes.NewBufferString(t.rawBody)

	req, err := http.NewRequest(t.method, t.host+"/"+t.path, body)

	for _, header := range t.headers {
		req.Header.Set(header.Field, header.Value)
	}

	if err != nil {
		fmt.Println("something went wrong !")
	}

	return req
}

type loader struct {
	n int
	c int

	tasks []*Task
	ws    *websocket.Conn
}

func New(n, c int, ws *websocket.Conn) *loader {
	tasks := make([]*Task, 0)
	return &loader{n, c, tasks, ws}
}

func (l *loader) AddTasks(method string, host string, path string, headers []models.Header, rawBody string) {
	t := &Task{method: method,
		host:    host,
		path:    path,
		headers: headers,
		rawBody: rawBody,
	}

	l.tasks = append(l.tasks, t)
}

func (l *loader) spawnClient(wg *sync.WaitGroup, queue chan []*Task) {
	// this is one out of c clients
	client := &http.Client{}

	// if we're not busy, get a tasks list from the queue channel
	for tasks := range queue {
		for i, task := range tasks {
			req := task.Request()
			start := time.Now()
			// Send the request
			res, err := client.Do(req)
			if err != nil {
				if err := l.ws.WriteJSON(feedback{Task: i, StatusCode: 0, Duration: time.Since(start)}); err != nil {
					fmt.Println("Couldn't write to the socket")
				}
			} else {
				if err := l.ws.WriteJSON(feedback{Task: i, StatusCode: res.StatusCode, Duration: time.Since(start)}); err != nil {
					fmt.Println("Couldn't write to the socket")
				}

				res.Body.Close()
			}
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
	queue := make(chan []*Task, l.n)

	// Spawn all our clients
	for i := 0; i < l.c; i++ {
		go func() {
			l.spawnClient(&wg, queue)
		}()
	}

	// Populate our channel with tasks arrays
	for i := 0; i < l.n; i++ {
		queue <- l.tasks
	}

	// Close our channel
	close(queue)

	// unblock
	wg.Wait()

	// close the socket
	l.ws.Close()
}
