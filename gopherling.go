package main

import (
	"encoding/json"
	"fmt"
	"github.com/gophergala/gopherling/loader"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

var (
	database        *mgo.Database
	databaseSession *mgo.Session
)

type Task struct {
	Method string `bson:"method" json:"method"`
	Path   string `bson:"path" json:"path"`
}

type Test struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	BaseUrl     string        `bson:"base_url" json:"base_url"`
	Requests    int           `bson:"requests" json:"requests"`
	Concurrency int           `bson:"concurrency" json:"concurrency"`
	Tasks       []Task        `bson:"tasks" json:"tasks"`
}

func showTests(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	results := make([]Test, 0)
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

	var t Test

	err := decoder.Decode(&t)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	tests := database.C("tests")

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
	var t Test

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
	var t Test

	err := database.C("tests").Find(bson.M{"_id": bson.ObjectIdHex(ps.ByName("id"))}).One(&t)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	l := loader.New(1000, 100)

	l.AddTasks()
	go func() {
		l.Run()
	}()

	w.WriteHeader(200)
}

func main() {
	// Initialize our database connection
	databaseSession, err := mgo.Dial("127.0.0.1:27017")
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
	router.POST("/api/tests/:id/start", startTest)

	// Catch-all (angular app)
	router.NotFound = http.FileServer(http.Dir("static")).ServeHTTP

	// Start listening
	fmt.Println("Gopherling server started on port 9410")
	http.ListenAndServe(":9410", router)
}
