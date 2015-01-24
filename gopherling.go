package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
)

func showTests(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprintf(w, "I want to see all tests\n")
}

func addTest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprintf(w, "I want to add a test!\n")
}

func showTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "I want to see a specific test, (id: %s)!\n", ps.ByName("id"))
}

func updateTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "I want to update a test, (id: %s)!\n", ps.ByName("id"))
}

func deleteTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "I want to delete a test, (id: %s)!\n", ps.ByName("id"))
}

func startTest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "I want to start a test, (id: %s)!\n", ps.ByName("id"))
}

func main() {
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

		fmt.Println("Gopherling server started on port 9410")
    http.ListenAndServe(":9410", router)
}
