package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

var (
	database        *mgo.Database
	databaseSession *mgo.Session
)

type Test struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	BaseUrl     string        `bson:"base_url"`
	Requests    int           `bson:"requests"`
	Concurrency int           `bson:"concurrency"`
}

func showTests(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "I want to see all tests\n")
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
		w.WriteHeader(201)
	}
}

func showTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "I want to see a specific test, (id: %s)!\n", ps.ByName("id"))
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
	fmt.Fprintf(w, "I want to start a test, (id: %s)!\n", ps.ByName("id"))
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
