package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gophergala/gopherling/loader"
	"github.com/gophergala/gopherling/models"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

var (
	database        *mgo.Database
	databaseSession *mgo.Session
	upgrader        = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func showTests(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	results := make([]models.Test, 0)
	err := database.C("tests").Find(bson.M{}).All(&results)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	js, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
}

func addTest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var t models.Test

	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	tests := database.C("tests")

	// i'm in a hurry :p
	if len(t.Id) > 0 {
		err = database.C("tests").RemoveId(t.Id)
	}

	t.Id = bson.NewObjectId()
	err = tests.Insert(&t)

	if err != nil {
		w.WriteHeader(500)
	} else {
		js, err := json.Marshal(t)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(js)
	}
}

func showTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var t models.Test

	err := database.C("tests").Find(bson.M{"_id": bson.ObjectIdHex(ps.ByName("id"))}).One(&t)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	js, err := json.Marshal(t)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(js)
}

func updateTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "I want to update a test, (id: %s)!\n", ps.ByName("id"))
}

func deleteTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := database.C("tests").RemoveId(bson.ObjectIdHex(ps.ByName("id")))

	if err != nil {
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
}

func startTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var t models.Test

	err = database.C("tests").Find(bson.M{"_id": bson.ObjectIdHex(ps.ByName("id"))}).One(&t)
	if err != nil {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Not found")); err != nil {
			fmt.Println("An error occured: the requested test doesn't exist")
		}
		return
	}

	l := loader.New(t.Requests, t.Concurrency, conn)

	for _, task := range t.Tasks {
		l.AddTasks(task.Method, t.BaseUrl, task.Path, task.Headers, task.RawBody)
	}

	l.Run()
}

func main() {
	// Database infos
	var dbHost, dbPort string
	flag.StringVar(&dbHost, "dbHost", "127.0.0.1", "mongoDB host")
	flag.StringVar(&dbPort, "dbPort", "27017", "mongoDB port")
	flag.Parse()

	// Initialize our database connection
	databaseSession, err := mgo.Dial(dbHost + ":" + dbPort)
	if err != nil {
		panic("Couldn't connect to the database server")
	} else {
		fmt.Println("Successfully connected to the database server")
	}

	defer databaseSession.Close()

	// Select the appropriate database
	database = databaseSession.DB("gopherling")

	// Create a new router instance
	router := httprouter.New()

	// View all the tests
	router.GET("/api/tests", showTests)

	// Add a test
	router.POST("/api/tests", addTest)

	// View a single test
	router.GET("/api/tests/:id", showTest)

	// Update a test
	router.PUT("/api/tests/:id", updateTest)

	// Delete a test
	router.DELETE("/api/tests/:id", deleteTest)

	// Start a test
	router.GET("/api/tests/:id/start", startTest)

	// Catch-all (angular app)
	router.NotFound = http.FileServer(http.Dir("static")).ServeHTTP

	// Start listening
	fmt.Println("Gopherling server started on port 9410")
	http.ListenAndServe(":9410", router)
}
