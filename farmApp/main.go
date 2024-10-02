package main

import (
	_ "farmApp/docs" // Required for Swagger documentation
	"farmApp/pkg/handler"
	"farmApp/pkg/persistence"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"sync"
)

var once sync.Once

// @title Farm Customer API
// @version 1.0
// @description This is a simple API for managing farm customers.
// @host localhost:8080
// @BasePath /
func main() {
	once.Do(func() {
		persistence.CreateFarmDB()
	})

	r := mux.NewRouter()

	// Serve Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Serve static HTML page
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/", homePageHandler)

	// Define API routes
	r.HandleFunc("/customers", logRequest(handler.GetCustomers)).Methods("GET")
	r.HandleFunc("/customers/{id}", logRequest(handler.GetCustomer)).Methods("GET")
	r.HandleFunc("/customers", logRequest(handler.AddCustomer)).Methods("POST")
	r.HandleFunc("/customers/{id}", logRequest(handler.UpdateCustomer)).Methods("PUT")
	r.HandleFunc("/customers/{id}", logRequest(handler.DeleteCustomer)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Home page handler for the static HTML page
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

// logRequest is a middleware that logs the request and catches panics
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		log.Printf("Handling request: %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
