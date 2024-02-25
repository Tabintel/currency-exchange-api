// main.go

package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/tabintel/currency-exchange-api/handlers"
)

func main() {
	// Create a new instance of the Gorilla mux router
	router := mux.NewRouter()

	// Define the route for the exchange rate endpoint
	router.HandleFunc("/exchange-rate", handlers.ExchangeRateHandler).Methods("POST")

	// Start the HTTP server on port 8080
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
