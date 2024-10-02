package main

import (
	_ "farmApp/docs" // swagger generated docs
	"farmApp/farmApp/pkg/persistence"
	"farmApp/pkg/routeHandler"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title Farm Customer API
// @version 1.0
// @description This is a simple API for managing farm customers.
// @host localhost:8080
// @BasePath /

func main() {
	persistence.CreateFarmDB()

	r := mux.NewRouter()

	// Serve Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Serve static HTML page
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/", homePageHandler)

	// Define API routes
	r.HandleFunc("/customers", routeHandler.GetCustomers).Methods("GET")
	r.HandleFunc("/customers/{id}", routeHandler.GetCustomer).Methods("GET")
	r.HandleFunc("/customers", routeHandler.AddCustomer).Methods("POST")
	r.HandleFunc("/customers/{id}", routeHandler.UpdateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{id}", routeHandler.DeleteCustomer).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Home page handler for the static HTML page
func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}
